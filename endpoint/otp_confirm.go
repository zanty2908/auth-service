package endpoint

import (
	db "auth-service/db/gen"
	"auth-service/token"
	"auth-service/utils"
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

type OTPConfirmParams struct {
	Aud      token.Audience
	Platform token.Platform
	OTP      string
	Phone    string
}

func (s *Module) ConfirmOTP(c context.Context, request interface{}) (interface{}, error) {
	req := request.(*OTPConfirmParams)

	curTime := time.Now()
	otpAuthParams := &db.GetOTPAuthParams{Phone: req.Phone, ExpiresAt: curTime, Aud: string(req.Aud), Platform: string(req.Platform)}
	otpAuth, err := s.repo.GetOTPAuth(c, otpAuthParams)
	if err != nil {
		log.Err(err).Msg("OTP confirm: can not fount otp")
		return nil, NewErrorRes(http.StatusBadRequest, errors.New("otp_incorrect"))
	}

	isOtpCorrect := otpAuth.Otp == req.OTP
	if !isOtpCorrect {
		nextEnteredTime := otpAuth.EnteredTimes + 1
		if nextEnteredTime >= int16(s.config.Auth.OTP.EnteredIncorrectlyTimes) {
			// Delete this otp auth
			err = s.repo.DeleteOTPAuthByID(c, otpAuth.ID)
			if err != nil {
				log.Err(err).Msg("OTP confirm: delete otp auth when over entered time error")
				return nil, utils.ErrorFailed
			}
			return nil, NewErrorRes(http.StatusBadRequest, errors.New("otp_expired"))
		} else {
			updOtpParams := &db.UpdateEnteredTimeParams{
				ID:           otpAuth.ID,
				EnteredTimes: nextEnteredTime,
			}
			err = s.repo.UpdateEnteredTime(c, updOtpParams)
			if err != nil {
				log.Err(err).Msg("OTP confirm: update entered time of otp auth error")
				return nil, utils.ErrorFailed
			}
		}

		return nil, NewErrorRes(http.StatusBadRequest, errors.New("otp_incorrect"))
	}

	user, err := s.shouldCreateUser(c, otpAuth.ID, string(req.Aud), req.Phone, string(req.Platform))
	if err != nil {
		log.Err(err).Msg("OTP confirm: shouldCreateUser error")
		return nil, err
	}

	email := ""
	if user.Email != nil {
		email = *user.Email
	}

	accessToken, refreshToken, accessPayload, refreshPayload, err := s.tokenMaker.CreateTokenPair(
		token.Audience(req.Aud),
		token.Platform(req.Platform),
		user.ID,
		user.Role,
		email,
		user.Phone,
		make(map[string]interface{}),
	)
	if err != nil {
		log.Error().Err(err).Msg("OTP confirm: senerate token")
		return nil, err
	}

	err = s.repo.ExecTx(c, func(q *db.Queries) error {
		// Insert oauth token
		oauthTokenParams := &db.CreateOAuthTokenParams{
			TokenID:          accessPayload.Id,
			Platform:         string(req.Platform),
			Aud:              string(req.Aud),
			UserID:           user.ID,
			AccessToken:      accessToken,
			RefreshToken:     refreshToken,
			AccessExpiresAt:  accessPayload.ExpiresAt,
			RefreshExpiresAt: refreshPayload.ExpiresAt,
		}
		_, err = q.CreateOAuthToken(c, oauthTokenParams)
		if err != nil {
			log.Err(err).Msg("OTP confirm: insert oauth token error")
			return err
		}

		// Deactivate otp auth
		err = q.DeleteOTPAuthByID(c, otpAuth.ID)
		if err != nil {
			log.Err(err).Msg("OTP confirm: delete otp auth error")
			return err
		}

		return nil
	})
	if err != nil {
		log.Err(err).Msg("OTP confirm: tx error")
		return nil, err
	}

	res := &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserResponse: mapUserResponse(user),
	}

	return res, nil
}

func (s *Module) shouldCreateUser(c context.Context, otpAuthId int32, aud, phone, platform string) (*db.User, error) {
	getUserParams := &db.GetUserByPhoneParams{Phone: phone, Aud: aud}
	user, err := s.repo.GetUserByPhone(c, getUserParams)
	if err == pgx.ErrNoRows {
		// New user
		createUser := &db.CreateUserParams{
			ID:    uuid.NewString(),
			Phone: phone,
			Aud:   aud,
		}
		user, err = s.repo.CreateUserTx(c, otpAuthId, createUser)
		if err != nil {
			log.Err(err).Msg("OTP confirm: create user tx error")
			return nil, NewErrorRes(http.StatusInternalServerError, err)
		}
	} else if err != nil {
		log.Err(err).Msg("OTP confirm: get user by phone error")
		return nil, NewErrorRes(http.StatusInternalServerError, err)
	}

	if user.DeletedAt != nil {
		return nil, NewErrorRes(http.StatusLocked, utils.ErrorBlocked)
	}

	return user, nil
}

type TokenResponse struct {
	*UserResponse
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
