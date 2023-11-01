package transport

import (
	"auth-service/endpoint"
	"auth-service/proto/service"
	"auth-service/token"
	"auth-service/utils"
	"context"
	"encoding/json"
	"errors"
	"net/http"

	kitgrpc "github.com/go-kit/kit/transport/grpc"
	kithttp "github.com/go-kit/kit/transport/http"
)

// Http
type OTPSendRequest struct {
	Phone string `json:"phone"`
}

func (s *Module) HandleHttpOTPSend() *kithttp.Server {
	return s.newHttpEndpoint(
		s.ep.OTPSend,
		s.decodeOTPSendHttp,
	)
}

func (s *Module) decodeOTPSendHttp(_ context.Context, r *http.Request) (interface{}, error) {
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

	req := new(OTPSendRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return nil, err
	}

	phone, err := utils.ValidatePhone(req.Phone)
	if err != nil {
		return nil, endpoint.NewErrorRes(http.StatusBadRequest, utils.ErrorInvalidPhone)
	}

	return &endpoint.OTPSendParams{
		Aud:      reqAud,
		Platform: reqPlatform,
		Phone:    phone,
	}, nil
}

func (s *Module) HandleGRPCOTPSend() kitgrpc.Handler {
	return s.newGRPCEndpoint(
		s.ep.OTPSend,
		s.decodeOTPSendGRPC,
		s.encodeOTPSendGRPC,
	)
}

func (s *Module) decodeOTPSendGRPC(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*service.SendOTPReq)
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

	return &endpoint.OTPSendParams{
		Aud:      req.Aud,
		Platform: req.Platform,
		Phone:    req.Phone,
	}, nil
}

func (s *Module) encodeOTPSendGRPC(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(string)
	return &service.SimpleRes{Data: res}, nil
}
