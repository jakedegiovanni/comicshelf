package server

import (
	"log/slog"
	"net/http"
	"time"
)

func serverLogger() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			method := r.Method
			url := r.URL.String()

			slog.Info(url, slog.String("method", method))

			next.ServeHTTP(w, r) // todo log response code
		}
		return http.HandlerFunc(fn)
	}
}

func queryDate() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if !r.URL.Query().Has("date") {
				slog.Debug("no date found in query, setting and redirecting")

				query := r.URL.Query()
				query.Set("date", time.Now().Format(justTheDateFormat))
				r.URL.RawQuery = query.Encode()

				http.Redirect(w, r, r.URL.String(), http.StatusFound)
				return
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
