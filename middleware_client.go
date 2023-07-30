package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
)

var pathRegex = regexp.MustCompile(`^/v1/public.*$`)

type addBase struct {
	next http.RoundTripper
}

func (a *addBase) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Host = "gateway.marvel.com"
	req.URL.Host = req.Host
	req.URL.Scheme = "https"
	if !pathRegex.MatchString(req.URL.Path) {
		req.URL.Path = fmt.Sprintf("/v1/public%s", req.URL.Path)
	}
	log.Println("sending to", req.URL.String())
	return a.next.RoundTrip(req)
}
