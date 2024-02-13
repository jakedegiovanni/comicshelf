package marvel

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/jakedegiovanni/comicshelf"
	"github.com/jakedegiovanni/comicshelf/comicclient"
)

type dataWrapper[T any] struct {
	Code            interface{}      `json:"code"`
	Status          string           `json:"status"`
	Copyright       string           `json:"copyright"`
	AttributionText string           `json:"attributionText"`
	AttributionHTML string           `json:"attributionHTML"`
	Etag            string           `json:"etag"`
	Data            dataContainer[T] `json:"data"`
}

type dataContainer[T any] struct {
	Offset  int `json:"offset"`
	Limit   int `json:"limit"`
	Total   int `json:"total"`
	Count   int `json:"count"`
	Results []T `json:"results"`
}

type item struct {
	Name        string `json:"name"`
	ResourceURI string `json:"resourceURI"`
}

type collection struct {
	Items []item `json:"items"`
}

type uri struct {
	Type string `json:"type"`
	Url  string `json:"url"`
}

type date struct {
	Type string `json:"type"`
	Date string `json:"date"`
}

type thumbnail struct {
	Path      string `json:"path"`
	Extension string `json:"extension"`
}

type baseResult struct {
	Id          int       `json:"id"`
	Title       string    `json:"title"`
	ResourceURI string    `json:"resourceURI"`
	Urls        []uri     `json:"urls"`
	Modified    string    `json:"modified"`
	Thumbnail   thumbnail `json:"thumbnail"`
}

type series struct {
	baseResult
	Comics collection `json:"comics"`
}

type comic struct {
	baseResult
	Format      string `json:"format"`
	IssueNumber int    `json:"issueNumber"`
	Series      item   `json:"series"`
	Dates       []date `json:"dates"`
}

type Client struct {
	client      *http.Client
	comicCache  *Cache[dataWrapper[comic]]
	seriesCache *Cache[dataWrapper[series]]
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
		comicCache:  NewCache[dataWrapper[comic]](),
		seriesCache: NewCache[dataWrapper[series]](),
	}
}

func (c *Client) GetWeeklyComics(ctx context.Context, t time.Time) (comicshelf.Page[comicshelf.Comic], error) {
	if t.Weekday() == time.Sunday {
		// sunday is when the cut off is however new comics show up on the monday, the offset set here makes expectations match reality
		t = t.AddDate(0, 0, -1)
	}

	first, last := c.weekRange(c.marvelUnlimitedDate(t))
	endpoint := fmt.Sprintf("/comics?format=comic&formatType=comic&noVariants=true&dateRange=%s,%s&hasDigitalIssue=true&orderBy=issueNumber&limit=100", first.Format(c.cfg.DateLayout), last.Format(c.cfg.DateLayout))

	comics, err := request[comic](endpoint, c.comicCache, c.client, c.logger)
	if err != nil {
		return comicshelf.Page[comicshelf.Comic]{}, err
	}

	comics1 := comicshelf.Page[comicshelf.Comic]{
		Total:   comics.Data.Total,
		Limit:   comics.Data.Limit,
		Offset:  comics.Data.Offset,
		Count:   comics.Data.Count,
		Results: make([]comicshelf.Comic, comics.Data.Count),
	}

	for _, comic := range comics.Data.Results {
		comic1 := comicshelf.Comic{
			Id:           comic.Id,
			Title:        comic.Title,
			Urls:         make([]comicshelf.Url, len(comic.Urls)),
			Thumbnail:    "",
			Format:       comic.Format,
			IssuerNumber: comic.IssueNumber,
			Dates:        make([]comicshelf.Date, len(comic.Dates)),
			Attribution:  comics.AttributionText,
		}

		for _, uri := range comic.Urls {
			comic1.Urls = append(comic1.Urls, comicshelf.Url{
				Type: uri.Type,
				Url:  uri.Url,
			})
		}

		for _, date := range comic.Dates {
			comic1.Dates = append(comic1.Dates, comicshelf.Date{
				Type: date.Type,
				Date: date.Date,
			})
		}

		comics1.Results = append(comics1.Results, comic1)
	}

	return comics1, nil
}

func (c *Client) GetComic(ctx context.Context, id int) (comicshelf.Comic, error) {
	endpoint := fmt.Sprintf("/comics/%d", id)
	comic, err := request[comic](endpoint, c.comicCache, c.client, c.logger)
	if err != nil {
		return comicshelf.Comic{}, err
	}

	return comicshelf.Comic{}, nil
}

func (c *Client) GetComicsWithinSeries(ctx context.Context, id int) ([]comicshelf.Comic, error) {
	endpoint := fmt.Sprintf("/series/%d/comics?format=comic&formatType=comic&noVariants=true&hasDigitalIssue=true&orderBy=issueNumber&limit=100", id)
	comics, err := request[comic](endpoint, c.comicCache, c.client, c.logger)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (c *Client) GetSeries(ctx context.Context, id int) (comicshelf.Series, error) {
	endpoint := fmt.Sprintf("/series/%d", id)
	series, err := request[series](endpoint, c.seriesCache, c.client, c.logger)
	if err != nil {
		return comicshelf.Series{}, err
	}

	return comicshelf.Series{}, nil
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

func request[T any](endpoint string, cache *Cache[dataWrapper[T]], client *http.Client, logger *slog.Logger) (*dataWrapper[T], error) {
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
			return &data, nil
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

	var d dataWrapper[T]
	err := json.NewDecoder(resp.Body).Decode(&d)
	if err != nil {
		return nil, fmt.Errorf("could not decode data wrapper: %w", err)
	}

	cache.Put(endpoint, d)
	return &d, nil
}
