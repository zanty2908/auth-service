package endpoint

import (
	db "auth-service/db/gen"
	"auth-service/utils"
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

type OTPSendParams struct {
	Aud, Phone, Platform string
}

func (s *Module) OTPSend(c context.Context, request interface{}) (interface{}, error) {
	req := request.(*OTPSendParams)

	curTime := time.Now()

	otpAuthParams := &db.GetOTPAuthParams{Phone: req.Phone, ExpiresAt: curTime, Aud: req.Aud, Platform: req.Platform}
	otpAuth, err := s.repo.GetOTPAuth(c, otpAuthParams)

	if err != nil && err != pgx.ErrNoRows {
		log.Err(err).Msg("OTP send: get otp auth error")
		return nil, err
	}

	// Has otp auth valid -> check resend time
	if err == nil && otpAuth != nil {
		if otpAuth.ResendAt.After(curTime) {
			return nil, NewErrorRes(http.StatusBadRequest, errors.New("otp_not_resend_yet"))
		}

		// Deactivate otp valid by phone
		delOTPAuthParams := &db.DeleteOTPAuthByPhoneParams{
			Phone:     req.Phone,
			ExpiresAt: curTime,
			Aud:       req.Aud,
			Platform:  req.Platform,
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
		Platform:  req.Platform,
		Aud:       req.Aud,
		Phone:     req.Phone,
		Otp:       otp,
		ResendAt:  curTime.Add(s.config.Auth.OTP.ResendDuration),
		ExpiresAt: curTime.Add(s.config.Auth.OTP.ExpiresDuration),
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
