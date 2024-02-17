package marvel

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jakedegiovanni/comicshelf"
	"github.com/jakedegiovanni/comicshelf/comicclient"
	"golang.org/x/sync/errgroup"
)

var _ comicshelf.ComicService = (*Client)(nil)
var _ comicshelf.SeriesService = (*Client)(nil)

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
		client: comicclient.NewClient(&cfg.Client, comicclient.MiddlewareChain(
			comicclient.AddBaseMiddleware(logger, cfg.Client.BaseURL), // todo would prefer this to be managed by comicclient since it comes from its config
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

	marvelComics, err := request[comic](endpoint, c.comicCache, c.client, c.logger)
	if err != nil {
		return comicshelf.Page[comicshelf.Comic]{}, err
	}

	comics := transformPage[comic, comicshelf.Comic](marvelComics.Data)

	for _, comic := range marvelComics.Data.Results {
		com, err := transformComic(comic, marvelComics.AttributionText)
		if err != nil {
			return comicshelf.Page[comicshelf.Comic]{}, err
		}
		comics.Results = append(comics.Results, com)
	}

	return comics, nil
}

func (c *Client) GetComic(ctx context.Context, id int) (comicshelf.Comic, error) {
	endpoint := fmt.Sprintf("/comics/%d", id)
	marvelComic, err := request[comic](endpoint, c.comicCache, c.client, c.logger)
	if err != nil {
		return comicshelf.Comic{}, err
	}

	if marvelComic.Data.Count == 0 {
		return comicshelf.Comic{}, fmt.Errorf("could not find comics for id: %d", id)
	}

	return transformComic(marvelComic.Data.Results[0], marvelComic.AttributionText)
}

func (c *Client) GetComicsWithinSeries(ctx context.Context, id int) ([]comicshelf.Comic, error) {
	endpoint := fmt.Sprintf("/series/%d/comics?format=comic&formatType=comic&noVariants=true&hasDigitalIssue=true&orderBy=issueNumber&limit=100", id)
	marvelComics, err := request[comic](endpoint, c.comicCache, c.client, c.logger)
	if err != nil {
		return nil, err
	}

	comics := make([]comicshelf.Comic, 0, marvelComics.Data.Count)
	for _, comic := range marvelComics.Data.Results {
		com, err := transformComic(comic, marvelComics.AttributionText)
		if err != nil {
			return []comicshelf.Comic{}, nil
		}
		comics = append(comics, com)
	}

	return comics, nil
}

func (c *Client) GetSeries(ctx context.Context, id int) (comicshelf.Series, error) {
	endpoint := fmt.Sprintf("/series/%d", id)
	series, err := request[series](endpoint, c.seriesCache, c.client, c.logger)
	if err != nil {
		return comicshelf.Series{}, err
	}

	if series.Data.Count == 0 {
		return comicshelf.Series{}, fmt.Errorf("could not find series with id: %d", id)
	}

	return transformSeries(ctx, series.Data.Results[0], series.AttributionText, c.GetComic)
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

func transformPage[C, P any](data dataContainer[C]) comicshelf.Page[P] {
	return comicshelf.Page[P]{
		Total:   data.Total,
		Limit:   data.Limit,
		Offset:  data.Offset,
		Count:   data.Count,
		Results: make([]P, 0, data.Count),
	}
}

func transformSeries(ctx context.Context, series series, attribution string, getComic func(context.Context, int) (comicshelf.Comic, error)) (comicshelf.Series, error) {
	s := comicshelf.Series{
		Id:        series.Id,
		Title:     series.Title,
		Urls:      make([]comicshelf.Url, 0, len(series.Urls)),
		Thumbnail: fmt.Sprintf("%s/portrait_uncanny.%s", series.Thumbnail.Path, series.Thumbnail.Extension),
		Comics:    make([]comicshelf.Comic, 0, len(series.Comics.Items)),
	}

	for _, uri := range series.Urls {
		s.Urls = append(s.Urls, transformUrl(uri))
	}

	g := new(errgroup.Group)
	g.SetLimit(len(series.Comics.Items))
	for i, comic := range series.Comics.Items {
		i, comic := i, comic // https://golang.org/doc/faq#closures_and_goroutines
		g.Go(func() error {
			id, err := extractId(comic.ResourceURI)
			if err != nil {
				return err
			}

			c, err := getComic(ctx, id)
			if err != nil {
				return err
			}

			s.Comics[i] = c
			return nil
		})
	}

	err := g.Wait()
	if err != nil {
		return comicshelf.Series{}, err
	}

	return s, nil
}

func transformComic(comic comic, attribution string) (comicshelf.Comic, error) {
	c := comicshelf.Comic{
		Id:           comic.Id,
		Title:        comic.Title,
		Urls:         make([]comicshelf.Url, 0, len(comic.Urls)),
		Thumbnail:    fmt.Sprintf("%s/portrait_uncanny.%s", comic.Thumbnail.Path, comic.Thumbnail.Extension),
		Format:       comic.Format,
		IssuerNumber: comic.IssueNumber,
		Attribution:  attribution,
	}

	seriesId, err := extractId(comic.Series.ResourceURI)
	if err != nil {
		return comicshelf.Comic{}, err
	}
	c.SeriesId = seriesId

	for _, uri := range comic.Urls {
		c.Urls = append(c.Urls, transformUrl(uri))
	}

	for _, date := range comic.Dates {
		if !strings.EqualFold(strings.ToLower(date.Type), "onsaledate") {
			continue
		}

		c.OnSaleDate = date.Date
	}

	return c, nil
}

func transformUrl(u uri) comicshelf.Url {
	return comicshelf.Url{
		Type: u.Type,
		Url:  u.Url,
	}
}

func extractId(s string) (int, error) {
	r := regexp.MustCompile(`/([0-9]+)/?`)
	matches := r.FindStringSubmatch(s)
	for i := len(matches) - 1; i >= 0; i-- {
		id, err := strconv.Atoi(matches[i])
		if err != nil {
			continue
		}

		return id, nil
	}

	return 0, fmt.Errorf("could not extract a valid id from: %s", s)
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
