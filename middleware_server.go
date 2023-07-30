package main

import (
	"log"
	"net/http"
	"runtime/debug"
)

type Middleware func(next http.HandlerFunc) http.HandlerFunc

func MiddlewareChain(chain ...Middleware) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		if len(chain) == 0 {
			return next
		}

		wrapped := next
		for i := len(chain) - 1; i > 0; i-- {
			wrapped = chain[i](wrapped)
		}

		return wrapped
	}
}

func RecoverHandler() Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			defer func(uri string) {
				if r := recover(); r != nil {
					log.Println("recovered handler", uri, r)
					debug.PrintStack()
				}
			}(r.URL.String())

			next.ServeHTTP(w, r)
		}
	}
}

func AllowedMethods(methods ...string) Middleware {
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
