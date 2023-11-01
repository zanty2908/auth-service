package endpoint

import (
	"context"
	"net/http"

	"github.com/rs/zerolog/log"
)

type ValidateTokenParams struct {
	Aud   string
	Token string
}

func (s *Module) ValidateToken(c context.Context, request interface{}) (interface{}, error) {
	req := request.(*ValidateTokenParams)

	// Validate token
	payload, err := s.tokenMaker.VerifyToken(req.Token)
	if err != nil {
		log.Err(err).Msg("Validate token: invalid")
		return nil, NewErrorRes(http.StatusForbidden, err)
	}

	if req.Aud != "" && req.Aud != payload.Audience {
		log.Error().Msg("Validate token: dont matching audience")
		return nil, NewErrorRes(http.StatusForbidden, err)
	}

	// Check if the token is blacklisted
	blacklisted, err := s.repo.Exists(c, payload.Id)
	if err != nil {
		// Handle Redis error
		log.Err(err).Msg("Validate token: redis error, can not check blacklist token")
	}

	if blacklisted {
		return nil, NewErrorRes(http.StatusForbidden, err)
	}

	return true, nil
}
