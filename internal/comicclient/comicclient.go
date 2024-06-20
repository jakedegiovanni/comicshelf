package comicclient

import (
	"net/http"
)

func New(cfg *Config, middleware ...Middleware) *http.Client {
	return &http.Client{
		Transport: MiddlewareChain(middleware...)(http.DefaultTransport),
		Timeout:   cfg.Timeout,
	}
}
