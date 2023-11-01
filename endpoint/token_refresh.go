package endpoint

import (
	db "auth-service/db/gen"
	"auth-service/token"
	"auth-service/utils"
	"context"
	"errors"
	"net/http"

	"github.com/rs/zerolog/log"
)

func (s *Module) RefreshToken(c context.Context, request interface{}) (interface{}, error) {
	req := request.(string)

	reqRefreshPayload, err := s.tokenMaker.VerifyRefreshToken(req)
	if err != nil {
		log.Err(err).Msg("Refresh token: verify token error")
		return nil, err
	}

	reqTokenID := reqRefreshPayload.Id
	oauthToken, err := s.repo.GetOAuthToken(c, reqTokenID)
	if err != nil {
		log.Err(err).Msg("Refresh token: get oauth token error")
		return nil, NewErrorRes(http.StatusForbidden, errors.New("token_expired"))
	}

	if req != oauthToken.RefreshToken {
		log.Error().Msg("Refresh token: refresh token do not machting")
		return nil, utils.ErrorFailed
	}

	accessToken, refreshToken, accessPayload, refreshPayload, err := s.tokenMaker.CreateTokenPair(
		token.Audience(reqRefreshPayload.Audience),
		token.Platform(reqRefreshPayload.Platform),
		reqRefreshPayload.Subject,
		reqRefreshPayload.Role,
		reqRefreshPayload.Email,
		reqRefreshPayload.Phone,
		reqRefreshPayload.Data,
	)
	if err != nil {
		log.Error().Err(err).Msg("Refresh Token: senerate token")
		return nil, err
	}

	s.repo.ExecTx(c, func(q *db.Queries) error {

		// Insert oauth token
		oauthTokenParams := &db.CreateOAuthTokenParams{
			TokenID:          accessPayload.Id,
			UserID:           oauthToken.UserID,
			AccessToken:      accessToken,
			RefreshToken:     refreshToken,
			AccessExpiresAt:  accessPayload.ExpiresAt,
			RefreshExpiresAt: refreshPayload.ExpiresAt,
		}
		_, err = q.CreateOAuthToken(c, oauthTokenParams)
		if err != nil {
			log.Err(err).Msg("Refresh Token: insert oauth token error")
			return err
		}

		// Deleted old token
		err = q.DeleteOAuthToken(c, reqTokenID)
		if err != nil {
			log.Err(err).Msg("Refresh Token: delete old oauth token error")
			return err
		}

		return nil
	})

	if err != nil {
		log.Err(err).Msg("Refresh Token: tx error")
		return nil, err
	}

	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
