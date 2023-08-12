package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"regexp"
)

var pathRegex = regexp.MustCompile(`^/v1/public.*$`)

type ClientMiddleware func(next http.RoundTripper) http.RoundTripper

func ClientMiddlewareChain(chain ...ClientMiddleware) ClientMiddleware {
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

type addBase struct {
	next   http.RoundTripper
	logger *slog.Logger
}

func AddBase(logger *slog.Logger) ClientMiddleware {
	return func(next http.RoundTripper) http.RoundTripper {
		return &addBase{
			next:   next,
			logger: logger,
		}
	}
}

func (a *addBase) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Host = "gateway.marvel.com"
	req.URL.Host = req.Host
	req.URL.Scheme = "https"
	if !pathRegex.MatchString(req.URL.Path) {
		req.URL.Path = fmt.Sprintf("/v1/public%s", req.URL.Path)
	}
	a.logger.Debug("sending to", slog.String("url", req.URL.String()))
	return a.next.RoundTrip(req)
}
