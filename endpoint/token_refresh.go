package endpoint

import (
	db "auth-service/db/gen"
	"auth-service/utils"
	"context"
	"encoding/json"
	"errors"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/rs/zerolog/log"
)

func (s *Module) RefreshTokenEndpoint() *kithttp.Server {
	return s.newHttpEndpoint(
		s.refreshToken,
		s.decodeRefreshTokenRequest,
	)
}

type RefreshTokenParams struct {
	RefreshToken string `json:"refreshToken"`
}

func (s *Module) decodeRefreshTokenRequest(c context.Context, r *http.Request) (interface{}, error) {
	req := new(RefreshTokenParams)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return nil, err
	}

	if len(req.RefreshToken) == 0 {
		log.Error().Msg("Refresh token: bad request")
		return nil, utils.ErrorBadRequest
	}

	return req, nil
}

func (s *Module) refreshToken(c context.Context, request interface{}) (interface{}, error) {
	req := request.(*RefreshTokenParams)

	reqRefreshPayload, err := s.tokenMaker.VerifyRefreshToken(req.RefreshToken)
	if err != nil {
		log.Err(err).Msg("Refresh token: verify token error")
		return nil, err
	}

	reqTokenID := reqRefreshPayload.ID
	oauthToken, err := s.repo.GetOAuthToken(c, reqTokenID)
	if err != nil {
		log.Err(err).Msg("Refresh token: get oauth token error")
		return nil, newErrorRes(http.StatusForbidden, errors.New("token_expired"))
	}

	if req.RefreshToken != oauthToken.RefreshToken {
		log.Error().Msg("Refresh token: refresh token do not machting")
		return nil, utils.ErrorFailed
	}

	accessToken, refreshToken, accessPayload, refreshPayload, err := s.tokenMaker.CreateTokenPair(reqRefreshPayload.Data)
	if err != nil {
		log.Error().Err(err).Msg("Refresh Token: senerate token")
		return nil, err
	}

	s.repo.ExecTx(c, func(q *db.Queries) error {

		// Insert oauth token
		oauthTokenParams := &db.CreateOAuthTokenParams{
			TokenID:          accessPayload.ID,
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
