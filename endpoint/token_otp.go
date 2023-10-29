package endpoint

import (
	db "auth-service/db/gen"
	"auth-service/db/repo"
	"auth-service/utils"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

func (s *Module) TokenOTPEndpoint() *kithttp.Server {
	return s.newHttpEndpoint(
		s.tokenOTP,
		s.decodeTokenOTPRequest,
	)
}

type TokenOTPParams struct {
	OTP   string `json:"otp"`
	Phone string `json:"phone" example:"+84909000999"`
}

func (s *Module) decodeTokenOTPRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := new(TokenOTPParams)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return nil, err
	}

	_, err := utils.ValidatePhone(req.Phone)
	if err != nil {
		return nil, utils.ErrorInvalidPhone
	}

	return req, nil
}

func (s *Module) tokenOTP(c context.Context, request interface{}) (interface{}, error) {
	req := request.(*TokenOTPParams)

	curTime := time.Now()
	otpAuthParams := &db.GetOTPAuthParams{Phone: req.Phone, ExpiresAt: curTime}
	otpAuth, err := s.repo.GetOTPAuth(c, otpAuthParams)
	if err != nil {
		return nil, utils.ErrorNotFound
	}

	isOtpCorrect := otpAuth.Otp == req.OTP
	if !isOtpCorrect {
		return nil, errors.New("user_otp_incorrect")
	}

	user, err := s.repo.GetUserByPhone(c, req.Phone)

	if err == pgx.ErrNoRows {
		// New user
		createUser := &db.CreateUserParams{
			ID:     uuid.NewString(),
			Phone:  req.Phone,
			Status: int16(repo.DRAFT),
		}
		user, err = s.repo.CreateUserTx(c, otpAuth.ID, createUser)
		if err != nil {
			log.Err(err).Msg("Token OTP: create user tx error")
			return nil, err
		}
	} else if err != nil {
		log.Err(err).Msg("Token OTP: get user by phone error")
		return nil, errors.New("user_incorrect")
	}

	if user.DeletedAt != nil {
		return nil, utils.ErrorBlocked
	}

	tokenData := map[string]string{
		"userId": user.ID,
		"phone":  user.Phone,
	}

	accessToken, refreshToken, accessPayload, refreshPayload, err := s.tokenMaker.CreateTokenPair(tokenData)
	if err != nil {
		log.Error().Err(err).Msg("Token OTP: senerate token")
		return nil, err
	}

	err = s.repo.ExecTx(c, func(q *db.Queries) error {
		// Insert oauth token
		oauthTokenParams := &db.CreateOAuthTokenParams{
			TokenID:          accessPayload.ID,
			UserID:           &user.ID,
			AccessToken:      accessToken,
			RefreshToken:     refreshToken,
			AccessExpiresAt:  accessPayload.ExpiresAt,
			RefreshExpiresAt: refreshPayload.ExpiresAt,
		}
		_, err = q.CreateOAuthToken(c, oauthTokenParams)
		if err != nil {
			log.Err(err).Msg("Token OTP: insert oauth token error")
			return err
		}

		// Deactivate otp auth
		err = q.DeleteOTPAuthByID(c, otpAuth.ID)
		if err != nil {
			log.Err(err).Msg("Token OTP: delete otp auth error")
			return err
		}

		return nil
	})
	if err != nil {
		log.Err(err).Msg("Token OTP: tx error")
		return nil, err
	}

	res := TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	res.UserResponse = mapUserResponse(user)

	return res, nil
}

type TokenResponse struct {
	*UserResponse
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
