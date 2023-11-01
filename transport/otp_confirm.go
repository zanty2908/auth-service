package transport

import (
	endpoint "auth-service/endpoint"
	"auth-service/proto/service"
	"auth-service/token"
	"auth-service/utils"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	kitgrpc "github.com/go-kit/kit/transport/grpc"
	kithttp "github.com/go-kit/kit/transport/http"
)

// Http
type OTPConfirmRequest struct {
	OTP   string `json:"otp"`
	Phone string `json:"phone" example:"+84909000999"`
}

func (s *Module) HandleHttpOTPConfirm() *kithttp.Server {
	return s.newHttpEndpoint(
		s.ep.ConfirmOTP,
		s.decodeOTPConfirmHttp,
	)
}

func (s *Module) decodeOTPConfirmHttp(_ context.Context, r *http.Request) (interface{}, error) {
	reqAud := r.URL.Query().Get("aud")
	aud := token.ParseAudience(reqAud)
	if aud == nil {
		return nil, endpoint.NewErrorRes(http.StatusBadRequest, errors.New("audience_missing"))
	}

	reqPlatform := r.URL.Query().Get("platform")
	platform := token.ParsePlatform(reqPlatform)
	if platform == nil {
		return nil, endpoint.NewErrorRes(http.StatusBadRequest, errors.New("platform_missing"))
	}

	req := new(OTPConfirmRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return nil, err
	}

	req.OTP = strings.TrimSpace(req.OTP)
	if req.OTP == "" {
		return nil, endpoint.BadRequestError()
	}

	_, err := utils.ValidatePhone(req.Phone)
	if err != nil {
		return nil, utils.ErrorInvalidPhone
	}

	return &endpoint.OTPConfirmParams{
		Aud:      *aud,
		Platform: *platform,
		Phone:    req.Phone,
		OTP:      req.OTP,
	}, nil
}

func (s *Module) HandleGRPCOTPConfirm() kitgrpc.Handler {
	return s.newGRPCEndpoint(
		s.ep.ConfirmOTP,
		s.decodeOTPConfirmGRPC,
		s.encodeOTPConfirmGRPC,
	)
}

func (s *Module) decodeOTPConfirmGRPC(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*service.ConfirmOTPReq)
	if req == nil {
		return nil, endpoint.BadRequestError()
	}

	aud := token.ParseAudience(req.Aud)
	if aud == nil {
		return nil, endpoint.NewErrorRes(http.StatusBadRequest, errors.New("audience_missing"))
	}

	platform := token.ParsePlatform(req.Platform)
	if platform == nil {
		return nil, endpoint.NewErrorRes(http.StatusBadRequest, errors.New("platform_missing"))
	}

	_, err := utils.ValidatePhone(req.Phone)
	if err != nil {
		return nil, utils.ErrorInvalidPhone
	}

	otp := strings.TrimSpace(req.Otp)
	if otp == "" {
		return nil, endpoint.BadRequestError()
	}

	return &endpoint.OTPConfirmParams{
		Aud:      *aud,
		Platform: *platform,
		Phone:    req.Phone,
		OTP:      otp,
	}, nil
}

func (s *Module) encodeOTPConfirmGRPC(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(*endpoint.TokenResponse)
	var email string
	if res.Email != nil {
		email = *res.Email
	}
	return &service.ConfirmOTPRes{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
		User: &service.User{
			Id:       res.ID,
			Phone:    res.Phone,
			Email:    email,
			Country:  res.Country,
			FullName: res.FullName,
			Aud:      res.Aud,
			Role:     res.Role,
		},
	}, nil
}
