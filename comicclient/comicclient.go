package comicclient

import (
	"net/http"
)

func NewClient(cfg *Config, middleware ...Middleware) *http.Client {
	return &http.Client{
		Transport: ClientMiddlewareChain(middleware...)(http.DefaultTransport),
		Timeout:   cfg.Timeout,
	}
}

func ClientMiddlewareChain(chain ...Middleware) Middleware {
	return func(next http.RoundTripper) http.RoundTripper {
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
