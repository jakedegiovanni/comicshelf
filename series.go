package comicshelf

import "context"

type Series struct {
	Name      string  `json:"name"`
	Comics    []Comic `json:"comics"`
	Id        int     `json:"id"`
	Title     string  `json:"title"`
	Urls      []Url   `json:"urls"`
	Thumbnail string  `json:"thumbnail"`
}

type SeriesService interface {
	GetComicsWithinSeries(ctx context.Context, id int) ([]Comic, error)
	GetSeries(ctx context.Context, id int) (Series, error)
}
