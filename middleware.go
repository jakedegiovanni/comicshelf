package main

import (
	"fmt"
	"net/http"
)

type addBase struct {
	next http.RoundTripper
}

func (a *addBase) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Host = "gateway.marvel.com"
	req.URL.Host = req.Host
	req.URL.Scheme = "https"
	req.URL.Path = fmt.Sprintf("/v1/public%s", req.URL.Path)
	return a.next.RoundTrip(req)
}
