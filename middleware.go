package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type apiKeyMiddleWare struct {
	next http.RoundTripper
	pub  io.Reader
	priv io.Reader
}

func (a *apiKeyMiddleWare) RoundTrip(req *http.Request) (*http.Response, error) {
	pub, err := io.ReadAll(a.pub)
	if err != nil {
		log.Fatalln("pub read", err)
	}

	priv, err = io.ReadAll(a.priv)
	if err != nil {
		log.Fatalln("priv read", err)
	}

	ts := fmt.Sprintf("%d", time.Now().UTC().Unix())
	hash := md5.Sum([]byte(ts + string(priv) + string(pub)))
	query := req.URL.Query()
	query.Add("ts", ts)
	query.Add("hash", fmt.Sprintf("%x", hash))
	query.Add("apikey", string(pub))
	req.URL.RawQuery = query.Encode()
	return a.next.RoundTrip(req)
}

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
