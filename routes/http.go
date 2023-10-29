package routes

import (
	"net/http"

	"auth-service/endpoint"
	"auth-service/middleware"
	"auth-service/server"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/gorilla/mux"
)

func NewHttpRouter(sv *server.Server, ep *endpoint.Module) http.Handler {
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "language"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Use(middleware.Logging)
	r.Use(middleware.Language)

	r.HandleFunc("/healthz", func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusOK)
	})

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r.Post("/otp", ep.OTPSendEndpoint().ServeHTTP)
	r.Route("/token", func(r chi.Router) {
		r.Get("/validate", ep.ValidateTokenEndpoint().ServeHTTP)
		r.Post("/", ep.TokenOTPEndpoint().ServeHTTP)
		r.Post("/refresh", ep.RefreshTokenEndpoint().ServeHTTP)
	})

	r.Route("/user", func(r chi.Router) {
		r.Use(middleware.Authen(sv))
		r.Get("/", ep.GetUserEndpoint().ServeHTTP)
	})

	return r
}

func setupMuxRouter(sv *server.Server, ep *endpoint.Module) {
	r := mux.NewRouter()

	r.Use(middleware.Logging)
	r.Use(middleware.Language)

	r.HandleFunc("/healthz", func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusOK)
	})

	r.Handle("/otp", ep.OTPSendEndpoint()).Methods("GET")
	r.Handle("/token", ep.TokenOTPEndpoint()).Methods("POST")
	r.Handle("/token/refresh", ep.RefreshTokenEndpoint()).Methods("POST")

	user := r.PathPrefix("/user").Subrouter()
	user.Use(middleware.Authen(sv))
	user.Handle("", ep.GetUserEndpoint()).Methods("GET")

}
