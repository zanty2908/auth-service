package transport

import (
	db "auth-service/db/gen"
	"auth-service/endpoint"
	"auth-service/proto/service"
	"auth-service/utils"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	kitgrpc "github.com/go-kit/kit/transport/grpc"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/rs/zerolog/log"
)

// Http
type UpdateUserParams struct {
	FullName *string `json:"fullName"`
	Email    *string `json:"email"`
}

func (s *Module) HandleHttpUserUpdate() *kithttp.Server {
	return s.newHttpEndpoint(
		s.ep.UserUpdate,
		s.decodeUserUpdateHttp,
	)
}

func (s *Module) decodeUserUpdateHttp(_ context.Context, r *http.Request) (interface{}, error) {
	payload, err := s.tokenMaker.GetTokenPayload(r)
	if err != nil {
		return nil, endpoint.BadRequestError()
	}

	req := new(UpdateUserParams)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		log.Err(err).Msg("Update User: parse request body error")
		return nil, endpoint.BadRequestError()
	}

	if req.Email != nil && !utils.ValidateEmail(*req.Email) {
		return nil, endpoint.NewErrorRes(http.StatusBadRequest, errors.New("invalid_email"))
	}

	params := &db.UpdateUserParams{
		ID:       payload.Subject,
		FullName: req.FullName,
		Email:    req.Email,
	}

	return params, nil
}

// GRPC
func (s *Module) HandleGRPCUserUpdate() kitgrpc.Handler {
	return s.newGRPCEndpoint(
		s.ep.UserUpdate,
		s.decodeUserUpdateGRPC,
		s.encodeUserUpdateGRPC,
	)
}

func (s *Module) decodeUserUpdateGRPC(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*service.UpdateUserReq)
	if req == nil {
		return nil, endpoint.BadRequestError()
	}

	reqId := strings.TrimSpace(req.Id)
	if len(reqId) == 0 {
		return nil, endpoint.BadRequestError()
	}

	if len(req.Email) > 0 && !utils.ValidateEmail(req.Email) {
		return nil, endpoint.NewErrorRes(http.StatusBadRequest, errors.New("invalid_email"))
	}

	return &db.UpdateUserParams{
		ID:       reqId,
		FullName: &req.FullName,
		Email:    &req.Email,
	}, nil
}

func (s *Module) encodeUserUpdateGRPC(_ context.Context, grpcRes interface{}) (interface{}, error) {
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
