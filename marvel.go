package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"
)

const (
	monthOffset = -3
	layout      = "2006-01-02"
)

type apiKeyMiddleWare struct {
	next   http.RoundTripper
	pub    io.ReadSeeker
	priv   io.ReadSeeker
	logger *slog.Logger
}

func ApiKeyMiddleware(logger *slog.Logger, pub, priv io.ReadSeeker) ClientMiddleware {
	return func(next http.RoundTripper) http.RoundTripper {
		return &apiKeyMiddleWare{
			next:   next,
			logger: logger,
			pub:    pub,
			priv:   priv,
		}
	}
}

func (a *apiKeyMiddleWare) RoundTrip(req *http.Request) (*http.Response, error) {
	_, err := a.pub.Seek(0, io.SeekStart)
	if err != nil {
		a.logger.Error(err.Error())
		os.Exit(1) // todo - we shouldn't be doing this
	}

	pub, err := io.ReadAll(a.pub)
	if err != nil {
		a.logger.Error("pub read", slog.String("err", err.Error()))
		os.Exit(1) // todo - we shouldn't be doing this
	}

	_, err = a.priv.Seek(0, io.SeekStart)
	if err != nil {
		a.logger.Error(err.Error())
		os.Exit(1) // todo - we shouldn't be doing this
	}

	priv, err = io.ReadAll(a.priv)
	if err != nil {
		a.logger.Error("priv read", slog.String("err", err.Error()))
		os.Exit(1) // todo - we shouldn't be doing this
	}

	ts := fmt.Sprintf("%d", time.Now().UTC().Unix())
	hash := md5.Sum([]byte(ts + string(priv) + string(pub)))
	query := req.URL.Query()
	query.Add("ts", ts)
	query.Add("hash", fmt.Sprintf("%x", hash))
	query.Add("apikey", string(pub))
	req.URL.RawQuery = query.Encode()
	a.logger.Debug("api key middleware")
	return a.next.RoundTrip(req)
}

type DataWrapper[T any] struct {
	Code            interface{}      `json:"code"`
	Status          string           `json:"status"`
	Copyright       string           `json:"copyright"`
	AttributionText string           `json:"attributionText"`
	AttributionHTML string           `json:"attributionHTML"`
	Etag            string           `json:"etag"`
	Data            DataContainer[T] `json:"data"`
}

type DataContainer[T any] struct {
	Offset  int `json:"offset"`
	Limit   int `json:"limit"`
	Total   int `json:"total"`
	Count   int `json:"count"`
	Results []T `json:"results"`
}

type Item struct {
	Name        string `json:"name"`
	ResourceURI string `json:"resourceURI"`
}

type Collection struct {
	Items []Item `json:"items"`
}

type Url struct {
	Type string `json:"type"`
	Url  string `json:"url"`
}

type Thumbnail struct {
	Path      string `json:"path"`
	Extension string `json:"extension"`
}

type BaseResult struct {
	Id          int       `json:"id"`
	Title       string    `json:"title"`
	ResourceURI string    `json:"resourceURI"`
	Urls        []Url     `json:"urls"`
	Modified    string    `json:"modified"`
	Thumbnail   Thumbnail `json:"thumbnail"`
}

type MarvelSeries struct {
	BaseResult
	Comics Collection `json:"comics"`
}

type MarvelComic struct {
	BaseResult
	Format      string `json:"format"`
	IssueNumber int    `json:"issueNumber"`
	Series      Item   `json:"series"`
}

type MarvelClient struct {
	client    *http.Client
	etagCache map[string]interface{}
	mu        *sync.Mutex
	logger    *slog.Logger
}

func NewMarvelClient(logger *slog.Logger) *MarvelClient {
	chain := ClientMiddlewareChain(
		AddBase(logger),
		ApiKeyMiddleware(logger, Pub, Priv),
	)

	return &MarvelClient{
		client: &http.Client{
			Timeout:   20 * time.Second,
			Transport: chain(http.DefaultTransport),
		},
		etagCache: make(map[string]interface{}),
		mu:        &sync.Mutex{},
		logger:    logger,
	}
}

func (m *MarvelClient) GetWeeklyComics(t time.Time) (*DataWrapper[MarvelComic], error) {
	if t.Weekday() == time.Sunday {
		// sunday is when the cut off is however new comics show up on the monday, the offset set here makes expectations match reality
		t = t.AddDate(0, 0, -1)
	}

	first, last := m.weekRange(t.AddDate(0, monthOffset, 0))
	endpoint := fmt.Sprintf("/comics?format=comic&formatType=comic&noVariants=true&dateRange=%s,%s&hasDigitalIssue=true&orderBy=issueNumber&limit=100", first.Format(layout), last.Format(layout))

	return request[MarvelComic](endpoint, m.mu, m.etagCache, m.client, m.logger)
}

func (m *MarvelClient) GetComic(endpoint string) (*DataWrapper[MarvelComic], error) {
	return request[MarvelComic](endpoint, m.mu, m.etagCache, m.client, m.logger)
}

func (m *MarvelClient) GetSeries(endpoint string) (*DataWrapper[MarvelSeries], error) {
	return request[MarvelSeries](endpoint, m.mu, m.etagCache, m.client, m.logger)
}

func (m *MarvelClient) weekRange(t time.Time) (time.Time, time.Time) {
	for t.Weekday() != time.Sunday {
		t = t.AddDate(0, 0, -1)
	}

	t = t.AddDate(0, 0, -7)

	return t, t.AddDate(0, 0, 6)
}

func request[T any](
	endpoint string,
	mu *sync.Mutex,
	etagCache map[string]interface{},
	client *http.Client,
	logger *slog.Logger,
) (*DataWrapper[T], error) {
	if resp, ok := cacheRead(endpoint, etagCache, mu); ok {
		resp, ok := resp.(*DataWrapper[T])
		if ok {
			req, err := http.NewRequest(http.MethodGet, endpoint, nil)
			if err != nil {
				return nil, err
			}

			req.Header.Set("If-None-Match", resp.Etag)

			r, err := client.Do(req)
			if err != nil {
				return nil, err
			}
			defer r.Body.Close()

			if r.StatusCode == http.StatusNotModified {
				logger.Debug("not modified, using cached response")
				return resp, nil
			}
		}
	}

	logger.Debug("item modified or not present in cache")

	resp, err := client.Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var d DataWrapper[T]
	err = json.NewDecoder(resp.Body).Decode(&d)
	if err != nil {
		return nil, err
	}

	cacheUpdate(endpoint, &d, etagCache, mu)
	return &d, nil
}

func cacheRead(key string, cache map[string]interface{}, mu *sync.Mutex) (interface{}, bool) {
	mu.Lock()
	defer mu.Unlock()
	i, ok := cache[key]
	return i, ok
}

func cacheUpdate(key string, val interface{}, cache map[string]interface{}, mu *sync.Mutex) {
	mu.Lock()
	defer mu.Unlock()
	cache[key] = val
}
