package middleware

import (
	"auth-service/endpoint"
	"auth-service/server"
	"encoding/json"
	"fmt"
	"net/http"
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

			for k, v := range payload.Data {
				r.Header.Add(k, fmt.Sprint(v))
			}

			next.ServeHTTP(w, r)
		})
	}
}
