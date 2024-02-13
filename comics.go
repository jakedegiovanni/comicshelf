package comicshelf

import (
	"context"
	"time"
)

type Comic struct {
	Name         string `json:"name"`
	Id           int    `json:"id"`
	Title        string `json:"title"`
	Urls         []Url  `json:"urls"`
	Thumbnail    string `json:"thumbnail"`
	Format       string `json:"format"`
	IssuerNumber int    `json:"issuer_number"`
	Dates        []Date `json:"dates"`
}

type ComicService interface {
	GetWeeklyComics(ctx context.Context, t time.Time) (Page[Comic], error)
	GetComic(ctx context.Context, id int) (Comic, error)
}
