package comicshelf

import (
	"context"
	"time"
)

type Comic struct {
	Id           int    `json:"id"`
	Title        string `json:"title"`
	Urls         []Url  `json:"urls"`
	Thumbnail    string `json:"thumbnail"`
	Format       string `json:"format"`
	IssuerNumber int    `json:"issuer_number"`
	OnSaleDate   string `json:"on_sale_date"`
	Attribution  string `json:"attribution"`
	SeriesId     int    `json:"series_id"`
}

type ComicService interface {
	GetWeeklyComics(ctx context.Context, t time.Time) (Page[Comic], error)
	GetComic(ctx context.Context, id int) (Comic, error)
}
