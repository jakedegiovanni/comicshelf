package comicshelf

type Page[T any] struct {
	Size    int `json:"size"`
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

type Date struct {
	Type string `json:"type"`
	Date string `json:"date"`
}
