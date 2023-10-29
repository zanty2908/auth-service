package endpoint

import (
	db "auth-service/db/gen"
	"auth-service/utils"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

func (s *Module) OTPSendEndpoint() *kithttp.Server {
	return s.newHttpEndpoint(
		s.otpSend,
		s.decodeOTPSendHttpReq,
	)
}

type OTPSendRequest struct {
	Phone string `json:"phone"`
}

func (s *Module) decodeOTPSendHttpReq(c context.Context, r *http.Request) (interface{}, error) {
	req := new(OTPSendRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return nil, err
	}

	phone, err := utils.ValidatePhone(req.Phone)
	if err != nil {
		return nil, utils.ErrorInvalidPhone
	}

	return phone, nil
}

func (s *Module) otpSend(c context.Context, request interface{}) (interface{}, error) {
	phone := request.(string)

	curTime := time.Now()

	otpAuthParams := &db.GetOTPAuthParams{Phone: phone, ExpiresAt: curTime}
	otpAuth, err := s.repo.GetOTPAuth(c, otpAuthParams)

	if err != nil && err != pgx.ErrNoRows {
		log.Err(err).Msg("OTP send: get otp auth error")
		return nil, err
	}

	// Has otp auth valid -> check resend time
	if err == nil && otpAuth != nil {
		if otpAuth.ResendAt.After(curTime) {
			return nil, errors.New("user_otp_not_resend_yet")
		}

		// Deactivate otp valid by phone
		delOTPAuthParams := &db.DeleteOTPAuthByPhoneParams{
			Phone:     phone,
			ExpiresAt: curTime,
		}
		err = s.repo.DeleteOTPAuthByPhone(c, delOTPAuthParams)
		if err != nil {
			log.Err(err).Msg("OTP send: delete old otp auth error")
			return nil, err
		}
	}

	otp, err := utils.GenerateOTP(4)
	if err != nil {
		log.Err(err).Msg("OTP send: generate otp error")
		return nil, utils.ErrorFailed
	}

	createOTPAuthParams := &db.CreateOTPAuthParams{
		Phone:     phone,
		Otp:       otp,
		ResendAt:  curTime.Add(s.config.Auth.OTP.ResendDuration),
		ExpiresAt: curTime.Add(s.config.Auth.OTP.ExpiresDuration),
		CreatedAt: curTime,
	}

	otpAuth, err = s.repo.CreateOTPAuth(c, createOTPAuthParams)
	if err != nil {
		log.Err(err).Msg("OTP send: create otp auth error")
		return nil, err
	}

	// Send OTP
	// err = channel.SendOTP(phone, otp)
	// if err != nil {
	// 	log.Error(localizer.MapError(err))
	// 	return response.NewError(c, http.StatusInternalServerError, localizer.MapError(err))
	// }

	return otp, nil
}
