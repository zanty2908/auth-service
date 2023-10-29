package middleware

import "net/http"

func Language(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		language := r.Header.Get("language")
		if language == "" {
			language = "en"
		}
		w.Header().Add("language", language)
		next.ServeHTTP(w, r)
	})
}
