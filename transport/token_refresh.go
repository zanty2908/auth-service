package transport

import (
	"auth-service/endpoint"
	"auth-service/proto/service"
	"auth-service/utils"
	"context"
	"encoding/json"
	"net/http"

	kitgrpc "github.com/go-kit/kit/transport/grpc"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/rs/zerolog/log"
)

// Http
type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}

func (s *Module) HandleHttpTokenRefresh() *kithttp.Server {
	return s.newHttpEndpoint(
		s.ep.RefreshToken,
		s.decodeTokenRefreshHttp,
	)
}

func (s *Module) decodeTokenRefreshHttp(_ context.Context, r *http.Request) (interface{}, error) {
	req := new(RefreshTokenRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return nil, err
	}

	if len(req.RefreshToken) == 0 {
		log.Error().Msg("Refresh token: bad request")
		return nil, utils.ErrorBadRequest
	}

	return req.RefreshToken, nil
}

func (s *Module) HandleGRPCTokenRefresh() kitgrpc.Handler {
	return s.newGRPCEndpoint(
		s.ep.RefreshToken,
		s.decodeTokenRefreshGRPC,
		s.encodeTokenRefreshGRPC,
	)
}

func (s *Module) decodeTokenRefreshGRPC(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*service.RefreshTokenReq)
	if len(req.Token) == 0 {
		log.Error().Msg("GRPC Refresh token: bad request")
		return nil, endpoint.BadRequestError()
	}

	return req.Token, nil
}

func (s *Module) encodeTokenRefreshGRPC(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(*endpoint.TokenResponse)
	return &service.RefreshTokenRes{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	}, nil
}
