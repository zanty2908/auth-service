package server

import (
	"auth-service/config"
	"auth-service/db/repo"
	"auth-service/language"
	"auth-service/token"

	"github.com/rs/zerolog/log"
)

type Server struct {
	Config         *config.Config
	Repo           repo.Repo
	TokenMaker     token.Maker
	MultiLocalizer *language.MultiLocalizer
}

func NewServer(config *config.Config, repo repo.Repo, multiLocalizer *language.MultiLocalizer) (*Server, error) {
	cfgToken := config.Auth.Token
	tokenMaker, err := token.NewJWTMaker(cfgToken.SecretKey, cfgToken.AccessTokenDuration, cfgToken.RefreshTokenDuration)
	if err != nil {
		log.Err(err).Msg("Create customer token maker error")
		return nil, err
	}

	return &Server{
		Config:         config,
		Repo:           repo,
		TokenMaker:     tokenMaker,
		MultiLocalizer: multiLocalizer,
	}, nil
}
