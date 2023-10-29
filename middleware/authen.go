package middleware

import (
	"auth-service/endpoint"
	"auth-service/server"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
)

func Authen(sv *server.Server) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			payload, err := sv.TokenMaker.GetTokenPayload(r)
			if err != nil {
				errMsg := err.Error()
				res := endpoint.Response{
					Meta: &endpoint.MetaResponse{
						Code:    http.StatusForbidden,
						Message: errMsg,
					},
				}
				b, _ := json.Marshal(res)
				w.WriteHeader(http.StatusForbidden)
				w.Write(b)
				return
			}

			payloadData := payload.Data

			if data, ok := payloadData.(map[string]interface{}); ok {
				for k, v := range data {
					r.Header.Add(k, fmt.Sprint(v))
				}
			} else {
				jsonData, err := json.Marshal(payloadData)
				if err != nil {
					log.Err(err).Msg("Auth Middleware: get data from payload error")
				}
				r.Header.Add("data", string(jsonData))
			}

			next.ServeHTTP(w, r)
		})
	}
}
