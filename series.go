package main

import (
	"html/template"
	"log"
	"net/http"
	"sort"
	"sync"
)

const SeriesEndpoint = "/marvel-unlimited/series"

type Series struct {
	tmpl   *template.Template
	client *MarvelClient
	db     *Db
}

func NewSeries(tmpl *template.Template, client *MarvelClient, db *Db) *Series {
	return &Series{
		tmpl:   tmpl,
		client: client,
		db:     db,
	}
}

func (s *Series) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.String())
	log.Println(r.URL.Query().Get("series"))

	resp, err := s.client.GetSeries(r.URL.Query().Get("series"))
	if err != nil {
		log.Fatalln("series", err)
	}

	var c chan []MarvelComic
	var wg sync.WaitGroup

	for _, result := range resp.Data.Results {
		c = make(chan []MarvelComic, len(result.Comics.Items))
		wg.Add(len(result.Comics.Items))
		for _, comic := range result.Comics.Items {
			go func(wg *sync.WaitGroup, c chan<- []MarvelComic, endpoint string) {
				defer wg.Done()

				w, err := s.client.GetComic(endpoint)
				if err != nil {
					log.Println("couldn't get comic", endpoint, err)
				}

				c <- w.Data.Results
			}(&wg, c, comic.ResourceURI)
		}
	}

	wg.Wait()
	close(c)

	result := &DataWrapper[MarvelComic]{
		Code:            resp.Code,
		Status:          resp.Status,
		Copyright:       resp.Copyright,
		AttributionText: resp.AttributionText,
		AttributionHTML: resp.AttributionHTML,
		Data: DataContainer[MarvelComic]{
			Results: make([]MarvelComic, 0),
		},
	}
	for comic := range c {
		result.Data.Results = append(result.Data.Results, comic...)
	}

	sort.SliceStable(result.Data.Results, func(i, j int) bool {
		return result.Data.Results[i].IssueNumber < result.Data.Results[j].IssueNumber
	})

	content := Content{
		Date:         r.URL.Query().Get("date"),
		PageEndpoint: SeriesEndpoint,
		Resp:         result,
	}

	err = s.tmpl.ExecuteTemplate(w, "index.html", content)
	if err != nil {
		log.Fatalln("exec", err)
	}
}
