package transport

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"auth-service/config"
	"auth-service/db/repo"
	"auth-service/endpoint"
	"auth-service/language"
	"auth-service/token"

	kitgrpc "github.com/go-kit/kit/transport/grpc"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Module struct {
	config         *config.Config
	repo           repo.Repo
	multiLocalizer *language.MultiLocalizer
	opts           []kithttp.ServerOption
	tokenMaker     token.Maker
	ep             *endpoint.Module
}

func NewModule(
	config *config.Config, repo repo.Repo,
	multiLocalizer *language.MultiLocalizer,
	tokenMaker token.Maker, ep *endpoint.Module,
) *Module {
	m := &Module{config: config, repo: repo, tokenMaker: tokenMaker, multiLocalizer: multiLocalizer, ep: ep}
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

// gRPC
func (s *Module) newGRPCEndpoint(
	handler func(context.Context, interface{}) (response interface{}, err error),
	decode func(context.Context, interface{}) (interface{}, error),
	encode func(context.Context, interface{}) (interface{}, error),
) *kitgrpc.Server {
	return kitgrpc.NewServer(
		handler,
		decode,
		encode,
	)
}

func (s *Module) encodeBaseGRPC(_ context.Context, grpcRes interface{}) (interface{}, error) {
	log.Info().Msgf(`Encode gRPC: %v`, grpcRes)
	return grpcRes, nil
}

func (s *Module) EncodeErrorGRPC(lang string, err error) error {
	log.Info().Msgf(`Encode error gRPC: %v`, err)

	localizer := s.multiLocalizer.Get(lang)

	var grpcErr *endpoint.ErrorResponse
	var grpcCode codes.Code
	if err == nil {
		grpcErr = endpoint.InternalServerError()
		grpcCode = codes.Internal
	} else if errRes, ok := err.(*endpoint.ErrorResponse); ok {
		code := errRes.StatusCode
		if code == 0 {
			code = endpoint.MappingErrorToCode(errRes.Err)
		}
		message := localizer.MapError(errRes.Err)
		grpcErr = endpoint.NewErrorRes(code, errors.New(message))
		grpcCode = codes.Unknown
	} else {
		grpcErr = endpoint.NewErrorRes(
			endpoint.MappingErrorToCode(err),
			errors.New(localizer.MapError(err)),
		)
		grpcCode = codes.Unknown
	}

	return status.Error(grpcCode, grpcErr.Error())
}
