package routes

import (
	"auth-service/proto/service"
	"auth-service/server"
	"auth-service/transport"
	"context"

	"github.com/rs/zerolog/log"
)

type gRPCServerImpl struct {
	transport *transport.Module
}

// SendOTP implements service.UserServiceServer.
func (s *gRPCServerImpl) SendOTP(c context.Context, req *service.SendOTPReq) (*service.SimpleRes, error) {
	_, res, err := s.transport.HandleGRPCOTPSend().ServeGRPC(c, req)
	if err != nil {
		log.Err(err).Msg("GRPC SendOTP error")
		return nil, s.transport.EncodeErrorGRPC(req.Language, err)
	}
	return res.(*service.SimpleRes), nil
}

// ConfirmOTP implements service.UserServiceServer.
func (s *gRPCServerImpl) ConfirmOTP(c context.Context, req *service.ConfirmOTPReq) (*service.ConfirmOTPRes, error) {
	_, res, err := s.transport.HandleGRPCOTPConfirm().ServeGRPC(c, req)
	if err != nil {
		log.Err(err).Msg("GRPC ConfirmOTP error")
		return nil, s.transport.EncodeErrorGRPC(req.Language, err)
	}
	return res.(*service.ConfirmOTPRes), nil
}

// GetUser implements service.UserServiceServer.
func (s *gRPCServerImpl) GetUser(c context.Context, req *service.GetUserReq) (*service.User, error) {
	_, res, err := s.transport.HandleGRPCUserGet().ServeGRPC(c, req)
	if err != nil {
		log.Err(err).Msg("GRPC GetUser error")
		return nil, s.transport.EncodeErrorGRPC(req.Language, err)
	}
	return res.(*service.User), nil
}

// RefreshToken implements service.UserServiceServer.
func (s *gRPCServerImpl) RefreshToken(c context.Context, req *service.RefreshTokenReq) (*service.RefreshTokenRes, error) {
	_, res, err := s.transport.HandleGRPCTokenRefresh().ServeGRPC(c, req)
	if err != nil {
		log.Err(err).Msg("GRPC RefreshToken error")
		return nil, s.transport.EncodeErrorGRPC(req.Language, err)
	}
	return res.(*service.RefreshTokenRes), nil
}

// UpdateUser implements service.UserServiceServer.
func (s *gRPCServerImpl) UpdateUser(c context.Context, req *service.UpdateUserReq) (*service.User, error) {
	_, res, err := s.transport.HandleGRPCUserUpdate().ServeGRPC(c, req)
	if err != nil {
		log.Err(err).Msg("GRPC UpdateUser error")
		return nil, s.transport.EncodeErrorGRPC(req.Language, err)
	}
	return res.(*service.User), nil
}

// ValidateToken implements service.UserServiceServer.
func (s *gRPCServerImpl) ValidateToken(c context.Context, req *service.ValidateTokenReq) (*service.SimpleRes, error) {
	_, res, err := s.transport.HandleGRPCTokenRefresh().ServeGRPC(c, req)
	if err != nil {
		log.Err(err).Msg("GRPC ValidateToken error")
		return nil, s.transport.EncodeErrorGRPC(req.Language, err)
	}
	return res.(*service.SimpleRes), nil
}

func NewGRPCRouter(sv *server.Server, transport *transport.Module) service.UserServiceServer {
	return &gRPCServerImpl{transport: transport}
}
