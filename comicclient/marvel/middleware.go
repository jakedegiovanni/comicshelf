package marvel

import (
	"crypto/md5"
	_ "embed"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/jakedegiovanni/comicshelf/comicclient"
)

// todo pub/priv will need to come from a secrets manager

//go:embed ..\..\pub.txt
var pub string

//go:embed ..\..\priv.txt
var priv string

func apiKeyMiddleware(logger *slog.Logger) comicclient.Middleware {
	return func(next http.RoundTripper) http.RoundTripper {
		return comicclient.MiddlewareFn(func(req *http.Request) (*http.Response, error) {
			ts := fmt.Sprintf("%d", time.Now().UTC().Unix())
			hash := md5.Sum([]byte(ts + priv + pub))
			query := req.URL.Query()
			query.Add("ts", ts)
			query.Add("hash", fmt.Sprintf("%x", hash))
			query.Add("apikey", pub)
			req.URL.RawQuery = query.Encode()

			logger.Debug("api key middleware")
			return next.RoundTrip(req)
		})
	}
}
