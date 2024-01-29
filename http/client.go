package http

import (
	"net/http"
	"time"
)

type ClientConfig struct {
	Timeout time.Duration
}

func NewClient(config ClientConfig, middleware ...func(http.RoundTripper) http.RoundTripper) *http.Client {
	rt := http.DefaultTransport.(*http.Transport).Clone()

	return &http.Client{
		Transport: ClientMiddlewareChain(middleware...)(rt),
		Timeout:   config.Timeout,
	}
}

func ClientMiddlewareChain(chain ...func(http.RoundTripper) http.RoundTripper) func(http.RoundTripper) http.RoundTripper {
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
