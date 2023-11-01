package endpoint

import (
	"auth-service/config"
	"auth-service/db/repo"
	"auth-service/token"
)

type Module struct {
	config     *config.Config
	repo       repo.Repo
	tokenMaker token.Maker
}

func NewEndpointModule(
	config *config.Config, repo repo.Repo, tokenMaker token.Maker) *Module {
	m := &Module{config: config, repo: repo, tokenMaker: tokenMaker}
	return m
}
