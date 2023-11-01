package endpoint

import (
	"context"
	"time"

	"auth-service/token"
	"auth-service/utils"

	"github.com/rs/zerolog/log"
)

func (s *Module) UserLogout(c context.Context, request interface{}) (interface{}, error) {
	payload := request.(*token.Payload)

	// Delete oauth token
	err := s.repo.DeleteOAuthToken(c, payload.Id)
	if err != nil {
		log.Err(err).Msg("Logout user: delete oauth token error")
		return nil, utils.ErrorFailed
	}

	// Check token has expired, if not, add token to blacklist with expires time
	diffTime := time.Now().Sub(payload.ExpiresAt)
	_, err = s.repo.Set(c, payload.Id, true, diffTime)
	if err != nil {
		log.Err(err).Msg("Logout user: add token to blacklist error")
	}

	return true, nil
}
