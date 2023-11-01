package transport

import (
	"auth-service/endpoint"
	"auth-service/proto/service"
	"auth-service/token"
	"auth-service/utils"
	"context"
	"errors"
	"net/http"
	"strings"

	kitgrpc "github.com/go-kit/kit/transport/grpc"
	kithttp "github.com/go-kit/kit/transport/http"
)

// Http
func (s *Module) HandleHttpTokenValidate() *kithttp.Server {
	return s.newHttpEndpoint(
		s.ep.ValidateToken,
		s.decodeTokenValidateHttp,
	)
}

func (s *Module) decodeTokenValidateHttp(_ context.Context, r *http.Request) (interface{}, error) {
	reqAud := r.URL.Query().Get("aud")
	aud := token.ParseAudience(reqAud)
	if aud == nil && reqAud != "" {
		return nil, endpoint.NewErrorRes(http.StatusBadRequest, errors.New("audience_missing"))
	}

	token := r.URL.Query().Get("token")
	token = strings.TrimSpace(token)
	if token == "" {
		return nil, utils.ErrorBadRequest
	}

	return &endpoint.ValidateTokenParams{
		Aud:   reqAud,
		Token: token,
	}, nil
}

func (s *Module) HandleGRPCTokenValidate() kitgrpc.Handler {
	return s.newGRPCEndpoint(
		s.ep.ValidateToken,
		s.decodeTokenValidateGRPC,
		s.encodeTokenValidateGRPC,
	)
}

func (s *Module) decodeTokenValidateGRPC(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*service.ValidateTokenReq)
	if req == nil {
		return nil, endpoint.BadRequestError()
	}

	aud := token.ParseAudience(req.Aud)
	if aud == nil {
		return nil, endpoint.NewErrorRes(http.StatusBadRequest, errors.New("audience_missing"))
	}

	token := strings.TrimSpace(req.Token)
	if token == "" {
		return nil, endpoint.BadRequestError()
	}

	return &endpoint.ValidateTokenParams{
		Aud:   req.Aud,
		Token: req.Token,
	}, nil
}

func (s *Module) encodeTokenValidateGRPC(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(string)
	return &service.SimpleRes{Data: res}, nil
}
