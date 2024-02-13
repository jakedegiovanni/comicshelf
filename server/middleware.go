package server

import (
	"log/slog"
	"net/http"
)

func serverLogger(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			method := r.Method
			url := r.URL.String()

			logger.Info(url, slog.String("method", method))

			next.ServeHTTP(w, r) // todo log response code
		}
		return http.HandlerFunc(fn)
	}
}
