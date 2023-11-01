package routes

import (
	"net/http"

	"auth-service/endpoint"
	"auth-service/middleware"
	"auth-service/server"
	"auth-service/transport"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func NewHttpRouter(sv *server.Server, transport *transport.Module, ep *endpoint.Module) http.Handler {
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "language"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Use(middleware.Logging)
	r.Use(middleware.Language)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r.Post("/otp/send", transport.HandleHttpOTPSend().ServeHTTP)
	r.Post("/otp/confirm", transport.HandleHttpOTPConfirm().ServeHTTP)
	r.Get("/token/validate", transport.HandleHttpTokenValidate().ServeHTTP)
	r.Post("/token/refresh", transport.HandleHttpTokenRefresh().ServeHTTP)

	r.Route("/user", func(r chi.Router) {
		r.Use(middleware.Authen(sv))
		r.Get("/", transport.HandleHttpUserGet().ServeHTTP)
		r.Put("/", transport.HandleHttpUserUpdate().ServeHTTP)
		r.Post("/logout", transport.HandleHttpUserLogout().ServeHTTP)
	})

	return r
}

// func NewMuxRouter(sv *server.Server, ep *endpoint.Module) http.Handler {
// 	r := mux.NewRouter()

// 	r.Use(cors.Handler(cors.Options{
// 		AllowedOrigins:   []string{"https://*", "http://*"},
// 		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
// 		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "language"},
// 		AllowCredentials: false,
// 		MaxAge:           300, // Maximum value not ignored by any of major browsers
// 	}))

// 	r.Use(middleware.Logging)
// 	r.Use(middleware.Language)

// 	r.HandleFunc("/health", func(rw http.ResponseWriter, req *http.Request) {
// 		rw.WriteHeader(http.StatusOK)
// 	})

// 	r.HandleFunc("/otp/send", ep.OTPSendEndpoint().ServeHTTP).Methods("POST")
// 	r.HandleFunc("/otp/confirm", ep.OTPConfirmEndpoint().ServeHTTP).Methods("POST")
// 	r.HandleFunc("/token/validate", ep.ValidateTokenEndpoint().ServeHTTP).Methods("GET")
// 	r.HandleFunc("/token/refresh", ep.RefreshTokenEndpoint().ServeHTTP).Methods("POST")

// 	user := r.PathPrefix("/user").Subrouter()
// 	user.Use(middleware.Authen(sv))
// 	user.HandleFunc("", ep.GetUserEndpoint().ServeHTTP).Methods("GET")
// 	user.HandleFunc("", ep.UpdateUserEndpoint().ServeHTTP).Methods("PUT")
// 	user.HandleFunc("/logout", ep.LogoutUserEndpoint().ServeHTTP).Methods("POST")

// 	return r
// }
