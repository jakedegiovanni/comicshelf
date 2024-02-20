package comicshelf

type Page[T any] struct {
	Limit   int `json:"limit"`
	Total   int `json:"total"`
	Count   int `json:"count"`
	Offset  int `json:"offset"`
	Results []T `json:"results"`
}

type Url struct {
	Type string `json:"type"`
	Url  string `json:"url"`
}
