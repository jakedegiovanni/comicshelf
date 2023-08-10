package main

import (
	"log/slog"
	"net/http"
	"runtime/debug"
)

type ServerMiddleware func(next http.HandlerFunc) http.HandlerFunc

func ServerMiddlewareChain(chain ...ServerMiddleware) ServerMiddleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		if len(chain) == 0 {
			return next
		}

		wrapped := next
		for i := len(chain) - 1; i >= 0; i-- {
			wrapped = chain[i](wrapped)
		}

		return wrapped
	}
}

func RecoverHandler(logger *slog.Logger) ServerMiddleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			defer func(uri string) {
				if r := recover(); r != nil {
					logger.Error("recovered handler", slog.String("url", uri), slog.Any("r", r))
					debug.PrintStack()
				}
			}(r.URL.String())

			next.ServeHTTP(w, r)
		}
	}
}

func AllowedMethods(methods ...string) ServerMiddleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			allowed := false
			for _, m := range methods {
				if m == r.Method {
					allowed = true
					break
				}
			}

			if !allowed {
				http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
				return
			}

			next.ServeHTTP(w, r)
		}
	}
}
