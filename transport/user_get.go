package transport

import (
	"auth-service/endpoint"
	"auth-service/proto/service"
	"context"
	"net/http"
	"strings"

	kitgrpc "github.com/go-kit/kit/transport/grpc"
	kithttp "github.com/go-kit/kit/transport/http"
)

// Http
func (s *Module) HandleHttpUserGet() *kithttp.Server {
	return s.newHttpEndpoint(
		s.ep.UserGet,
		s.decodeUserGetHttp,
	)
}

func (s *Module) decodeUserGetHttp(_ context.Context, r *http.Request) (interface{}, error) {
	payload, err := s.tokenMaker.GetTokenPayload(r)
	if err != nil {
		return nil, endpoint.BadRequestError()
	}
	return &endpoint.GetUserByIdParams{Id: payload.Subject}, nil
}

func (s *Module) HandleGRPCUserGet() kitgrpc.Handler {
	return s.newGRPCEndpoint(
		s.ep.UserGet,
		s.decodeUserGetGRPC,
		s.encodeUserGetGRPC,
	)
}

func (s *Module) decodeUserGetGRPC(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*service.GetUserReq)
	if req == nil {
		return nil, endpoint.BadRequestError()
	}

	reqId := strings.TrimSpace(req.Id)
	if len(reqId) == 0 {
		return nil, endpoint.BadRequestError()
	}

	return &endpoint.GetUserByIdParams{Id: reqId}, nil
}

func (s *Module) encodeUserGetGRPC(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(*endpoint.UserResponse)
	var email string
	if res.Email != nil {
		email = *res.Email
	}
	return &service.User{
		Id:       res.ID,
		Phone:    res.Phone,
		Email:    email,
		Country:  res.Country,
		FullName: res.FullName,
		Aud:      res.Aud,
		Role:     res.Role,
	}, nil
}
