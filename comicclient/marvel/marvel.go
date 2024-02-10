package marvel

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
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
	cfg         *Config
}

func New(cfg *Config, logger *slog.Logger) *Client {
	return &Client{
		client: comicclient.NewClient(nil, comicclient.MiddlewareChain(
			comicclient.AddBaseMiddleware(logger, &cfg.Client.BaseURL), // todo would prefer this to be managed by comicclient since it comes from its config
			apiKeyMiddleware(logger),
		)),
		logger:      logger,
		cfg:         cfg,
		comicCache:  NewCache[*DataWrapper[Comic]](),
		seriesCache: NewCache[*DataWrapper[Series]](),
	}
}

func (c *Client) GetWeeklyComics(t time.Time) (*DataWrapper[Comic], error) {
	if t.Weekday() == time.Sunday {
		// sunday is when the cut off is however new comics show up on the monday, the offset set here makes expectations match reality
		t = t.AddDate(0, 0, -1)
	}

	first, last := c.weekRange(c.marvelUnlimitedDate(t))
	endpoint := fmt.Sprintf("/comics?format=comic&formatType=comic&noVariants=true&dateRange=%s,%s&hasDigitalIssue=true&orderBy=issueNumber&limit=100", first.Format(c.cfg.DateLayout), last.Format(c.cfg.DateLayout))

	return request[Comic](endpoint, c.comicCache, c.client, c.logger)
}

func (c *Client) GetComic(endpoint string) (*DataWrapper[Comic], error) {
	return request[Comic](endpoint, c.comicCache, c.client, c.logger)
}

func (c *Client) GetComicsWithinSeries(seriesEndpoint string) (*DataWrapper[Comic], error) {
	uri, err := url.Parse(seriesEndpoint)
	if err != nil {
		return nil, fmt.Errorf("could not parse endpoint for series comics retrieval: %w", err)
	}

	uri.Path = fmt.Sprintf("%s/comics", uri.Path)
	query := uri.Query()
	query.Add("format", "comic")
	query.Add("formatType", "comic")
	query.Add("noVariants", "true")
	query.Add("hasDigitalIssue", "true")
	query.Add("orderBy", "issueNumber")
	query.Add("limit", "100")
	uri.RawQuery = query.Encode()

	return request[Comic](uri.String(), c.comicCache, c.client, c.logger)
}

func (c *Client) weekRange(t time.Time) (time.Time, time.Time) {
	for t.Weekday() != time.Sunday {
		t = t.AddDate(0, 0, -1)
	}

	t = t.AddDate(0, 0, -7)

	return t, t.AddDate(0, 0, 6)
}

func (c *Client) marvelUnlimitedDate(t time.Time) time.Time {
	return t.AddDate(0, c.cfg.ReleaseOffset, 0)
}

func request[T any](endpoint string, cache *Cache[*DataWrapper[T]], client *http.Client, logger *slog.Logger) (*DataWrapper[T], error) {
	var resp *http.Response

	if data, ok := cache.Get(endpoint); ok {
		req, err := http.NewRequest(http.MethodGet, endpoint, nil)
		if err != nil {
			return nil, fmt.Errorf("could not create request: %w", err)
		}

		req.Header.Set("If-None-Match", data.Etag)

		resp, err = client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("error whilst performing request: %w", err)
		}

		if resp.StatusCode == http.StatusNotModified {
			logger.Debug("not modified, using cached response")
			return data, nil
		}
	} else {
		logger.Debug("item not present in cache")

		var err error
		resp, err = client.Get(endpoint)
		if err != nil {
			return nil, fmt.Errorf("error whilst performing request: %w", err)
		}
	}

	defer resp.Body.Close()

	var d DataWrapper[T]
	err := json.NewDecoder(resp.Body).Decode(&d)
	if err != nil {
		return nil, fmt.Errorf("could not decode data wrapper: %w", err)
	}

	cache.Put(endpoint, &d)
	return &d, nil
}
