package marvel

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/jakedegiovanni/comicshelf/comicclient"
)

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

type Date struct {
	Type string `json:"type"`
	Date string `json:"date"`
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

type Series struct {
	BaseResult
	Comics Collection `json:"comics"`
}

type Comic struct {
	BaseResult
	Format      string `json:"format"`
	IssueNumber int    `json:"issueNumber"`
	Series      Item   `json:"series"`
	Dates       []Date `json:"dates"`
}

type Client struct {
	client      *http.Client
	comicCache  *Cache[*DataWrapper[Comic]]
	seriesCache *Cache[*DataWrapper[Series]]
	logger      *slog.Logger
}

func New(cfg *Config, logger *slog.Logger) *Client {
	return &Client{
		client: comicclient.NewClient(nil, comicclient.MiddlewareChain(
			comicclient.AddBaseMiddleware(logger, &cfg.Client.BaseURL), // todo would prefer this to be managed by comicclient since it comes from its config
			apiKeyMiddleware(logger),
		)),
		logger:      logger,
		comicCache:  NewCache[*DataWrapper[Comic]](),
		seriesCache: NewCache[*DataWrapper[Series]](),
	}
}

func (c *Client) GetWeeklyComics(t time.Time) (*DataWrapper[Comic], error) {
	return nil, errors.New("implement me")
}

func (c *Client) GetComic(endpoint string) (*DataWrapper[Comic], error) {
	var resp *http.Response

	if data, ok := c.comicCache.Get(endpoint); ok {
		req, err := http.NewRequest(http.MethodGet, endpoint, nil)
		if err != nil {
			return nil, fmt.Errorf("could not create request: %w", err)
		}

		req.Header.Set("If-None-Match", data.Etag)

		resp, err = c.client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("error whilst performing request: %w", err)
		}

		if resp.StatusCode == http.StatusNotModified {
			c.logger.Debug("not modified, using cached response")
			return data, nil
		}
	} else {
		c.logger.Debug("item not present in cache")

		var err error
		resp, err = c.client.Get(endpoint)
		if err != nil {
			return nil, fmt.Errorf("error whilst performing request: %w", err)
		}
	}

	defer resp.Body.Close()

	var d DataWrapper[Comic]
	err := json.NewDecoder(resp.Body).Decode(&d)
	if err != nil {
		return nil, fmt.Errorf("could not decode data wrapper: %w", err)
	}

	c.comicCache.Put(endpoint, &d)
	return &d, nil
}

func (c *Client) GetComicsWithinSeries(seriesEndpoint string) (*DataWrapper[Comic], error) {
	return nil, errors.New("implement me")
}
