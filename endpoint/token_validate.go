package endpoint

import (
	"auth-service/utils"
	"context"
	"net/http"
	"strings"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/rs/zerolog/log"
)

func (s *Module) ValidateTokenEndpoint() *kithttp.Server {
	return s.newHttpEndpoint(
		s.validateToken,
		s.decodeValidateTokenRequest,
	)
}

func (s *Module) decodeValidateTokenRequest(c context.Context, r *http.Request) (interface{}, error) {
	token := r.URL.Query().Get("token")
	token = strings.TrimSpace(token)
	if token == "" {
		return nil, utils.ErrorBadRequest
	}
	return token, nil
}

func (s *Module) validateToken(c context.Context, request interface{}) (interface{}, error) {
	token := request.(string)

	// Validate token
	payload, err := s.tokenMaker.VerifyToken(token)
	if err != nil {
		return nil, newErrorRes(http.StatusForbidden, err)
	}

	// Check if the token is blacklisted
	blacklisted, err := s.repo.Exists(c, payload.ID)
	if err != nil {
		// Handle Redis error
		log.Err(err).Msg("Validate token: redis error, can not check blacklist token")
	}

	if blacklisted {
		return nil, newErrorRes(http.StatusForbidden, err)
	}

	return true, nil
}
