package endpoint

import (
	"context"
	"encoding/json"
	"net/http"

	"auth-service/config"
	"auth-service/db/repo"
	"auth-service/language"
	"auth-service/token"

	kithttp "github.com/go-kit/kit/transport/http"
)

type Module struct {
	config         *config.Config
	repo           repo.Repo
	tokenMaker     token.Maker
	multiLocalizer *language.MultiLocalizer
	opts           []kithttp.ServerOption
}

func NewEndpointModule(config *config.Config, repo repo.Repo, multiLocalizer *language.MultiLocalizer, tokenMaker token.Maker) *Module {
	m := &Module{config: config, repo: repo, tokenMaker: tokenMaker, multiLocalizer: multiLocalizer}
	m.opts = []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(m.encodeError),
	}
	return m
}

func (s *Module) newHttpEndpoint(
	handler func(context.Context, interface{}) (response interface{}, err error),
	decode func(context.Context, *http.Request) (interface{}, error),
) *kithttp.Server {
	return kithttp.NewServer(
		handler,
		decode,
		s.encodeHttpResponse,
		s.opts...,
	)
}

func (s *Module) encodeHttpResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Authorization, content-type")

	language := w.Header().Get("language")
	localizer := s.multiLocalizer.Get(language)

	code, baseRes := mappingResponse(&localizer, response)
	w.WriteHeader(code)

	return json.NewEncoder(w).Encode(baseRes)
}

// encode errors from business-logic
func (s *Module) encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Authorization, content-type")

	language := w.Header().Get("language")
	localizer := s.multiLocalizer.Get(language)
	code, baseRes := mappingResponse(&localizer, err)
	w.WriteHeader(code)

	json.NewEncoder(w).Encode(baseRes)
}
