package token

import (
	"testing"
	"time"

	"auth-service/utils"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	random := utils.NewUtilRandom()
	maker, err := NewJWTMaker(random.RandomString(32), time.Minute, time.Hour)
	require.NoError(t, err)

	customerId := random.RandomString(12)
	phone := random.RandomPhone()

	issuedAt := time.Now()
	accessExpiresAt := issuedAt.Add(time.Minute)
	refreshExpiresAt := issuedAt.Add(time.Hour)

	data := map[string]interface{}{
		"cusId": customerId,
		"phone": phone,
	}

	accessToken, refreshToken, accessPayload, refreshPayload, err := maker.CreateTokenPair(
		USER_APP,
		ANDROID,
		random.RandomString(12),
		random.RandomString(10),
		random.RandomEmail(),
		random.RandomPhone(),
		data,
	)
	require.NoError(t, err)
	require.NotEmpty(t, accessToken)
	require.NotEmpty(t, refreshToken)
	require.NotEmpty(t, accessPayload)
	require.NotEmpty(t, refreshPayload)

	payload, err := maker.VerifyToken(accessToken)
	require.NoError(t, err)
	require.NotNil(t, payload)
	require.NotZero(t, payload.ExpiresAt)
	require.NotZero(t, payload.Id)
	require.Equal(t, data, payload.Data)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, accessExpiresAt, payload.ExpiresAt, time.Second)
	require.Equal(t, payload.Id, accessPayload.Id)

	rePayload, err := maker.VerifyRefreshToken(refreshToken)
	require.NoError(t, err)
	require.NotNil(t, rePayload)
	require.NotZero(t, rePayload.ExpiresAt)
	require.NotZero(t, rePayload.Id)
	require.Equal(t, data, rePayload.Data)
	require.WithinDuration(t, issuedAt, rePayload.IssuedAt, time.Second)
	require.WithinDuration(t, refreshExpiresAt, rePayload.ExpiresAt, time.Second)
	require.Equal(t, refreshPayload.Id, rePayload.Id)
}

func TestExpiredJWTToken(t *testing.T) {
	random := utils.NewUtilRandom()
	maker, err := NewJWTMaker(random.RandomString(32), -time.Minute, -time.Hour)
	require.NoError(t, err)

	accessToken, refreshToken, accessPayload, refreshPayload, err := maker.CreateTokenPair(
		USER_APP,
		ANDROID,
		random.RandomString(12),
		random.RandomString(10),
		random.RandomEmail(),
		random.RandomPhone(),
		make(map[string]interface{}, 0),
	)
	require.NoError(t, err)
	require.NotEmpty(t, accessToken)
	require.NotEmpty(t, accessPayload)
	require.NotEmpty(t, refreshToken)
	require.NotEmpty(t, refreshPayload)

	payload, err := maker.VerifyToken(accessToken)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)

	rePayload, err := maker.VerifyRefreshToken(refreshToken)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, rePayload)
}

func TestInvalidJWTTokenAlgNone(t *testing.T) {
	random := utils.NewUtilRandom()

	payload, err := NewPayload(
		USER_APP,
		ANDROID,
		random.RandomString(12),
		random.RandomString(10),
		random.RandomEmail(),
		random.RandomPhone(),
		make(map[string]interface{}, 0),
		time.Minute, nil,
	)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	maker, err := NewJWTMaker(random.RandomString(32), time.Minute, time.Hour)
	require.NoError(t, err)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}
