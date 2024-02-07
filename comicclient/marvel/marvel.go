package marvel

import (
	"log/slog"
	"net/http"

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
	comicCache  *Cache[Comic]
	seriesCache *Cache[Series]
	logger      *slog.Logger
}

func New(cfg *Config, logger *slog.Logger) *Client {
	return &Client{
		client: comicclient.NewClient(nil, comicclient.MiddlewareChain(
			comicclient.AddBaseMiddleware(logger, &cfg.Client.BaseURL), // todo would prefer this to be managed by comicclient since it comes from its config
			apiKeyMiddleware(logger),
		)),
		logger:      logger,
		comicCache:  NewCache[Comic](),
		seriesCache: NewCache[Series](),
	}
}
