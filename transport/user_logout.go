package transport

import (
	"auth-service/endpoint"
	"context"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
)

// Http
func (s *Module) HandleHttpUserLogout() *kithttp.Server {
	return s.newHttpEndpoint(
		s.ep.UserLogout,
		s.decodeUserLogoutHttp,
	)
}

func (s *Module) decodeUserLogoutHttp(_ context.Context, r *http.Request) (interface{}, error) {
	payload, err := s.tokenMaker.GetTokenPayload(r)
	if err != nil {
		return nil, endpoint.BadRequestError()
	}
	return payload, nil
}
