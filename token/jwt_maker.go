package token

import (
	"auth-service/utils"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/rs/zerolog/log"
)

const minSecretKeySize = 32

// JWTMaker is a JSON Web Token maker
type JWTMaker struct {
	accessKey       string
	accessDuration  time.Duration
	refreshKey      string
	refreshDuration time.Duration
}

// GetTokenPayload implements Maker
func (s *JWTMaker) GetTokenPayload(r *http.Request) (*Payload, error) {
	authorizationHeader := r.Header.Get("authorization")

	fields := strings.Fields(authorizationHeader)
	if len(fields) < 2 {
		err := errors.New("invalid authorization header format")
		return nil, err
	}

	accessToken := fields[1]
	payload, err := s.VerifyToken(accessToken)
	if err != nil {
		log.Err(err).Msg("GetTokenPayload error")
		return nil, err
	}

	return payload, nil
}

// NewJWTMaker creates a new JWTMaker
func NewJWTMaker(secretKey string, accessDuration, refreshDuration time.Duration) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}

	refreshSecretKey := utils.ReverseString(secretKey)
	refreshSecretKey = "rftk" + refreshSecretKey[4:]

	return &JWTMaker{
		accessKey:       secretKey,
		accessDuration:  accessDuration,
		refreshKey:      refreshSecretKey,
		refreshDuration: refreshDuration,
	}, nil
}

// CreateToken creates a new token for a specific username and duration
func (maker *JWTMaker) CreateTokenPair(aud Audience, platform Platform, userId, role, email, phone string, data map[string]interface{}) (string, string, *Payload, *Payload, error) {
	// Generate access token
	accessPayload, err := NewPayload(aud, platform, userId, role, email, phone, data, maker.accessDuration, nil)
	if err != nil {
		return "", "", accessPayload, nil, err
	}

	jwtAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessPayload)
	token, err := jwtAccessToken.SignedString([]byte(maker.accessKey))

	// Generate refresh token
	refreshPayload, err := NewPayload(aud, platform, userId, role, email, phone, data, maker.refreshDuration, &accessPayload.Id)
	if err != nil {
		return "", "", accessPayload, refreshPayload, err
	}

	jwtRefreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshPayload)
	refreshToken, err := jwtRefreshToken.SignedString([]byte(maker.refreshKey))

	return token, refreshToken, accessPayload, refreshPayload, err
}

// VerifyToken checks if the token is valid or not
func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(maker.accessKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}

func (maker *JWTMaker) VerifyRefreshToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(maker.refreshKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}
