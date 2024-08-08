package comicclient

import (
	"log/slog"
	"net/http"
	"net/url"
	"strings"
)

type Middleware func(http.RoundTripper) http.RoundTripper

type MiddlewareFn func(*http.Request) (*http.Response, error)

func (m MiddlewareFn) RoundTrip(req *http.Request) (*http.Response, error) {
	return m(req)
}

func MiddlewareChain(chain ...Middleware) Middleware {
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

func AddBaseMiddleware(url *url.URL) Middleware {
	return func(next http.RoundTripper) http.RoundTripper {
		return MiddlewareFn(func(req *http.Request) (*http.Response, error) {
			req.Host = url.Host
			req.URL.Host = url.Host
			req.URL.Scheme = url.Scheme

			if !strings.HasSuffix(req.URL.Path, url.Path) {
				u := *url
				req.URL.Path = u.JoinPath(req.URL.Path).Path
			}

			slog.Debug("sending to", slog.String("url", req.URL.String()))
			return next.RoundTrip(req)
		})
	}
}
