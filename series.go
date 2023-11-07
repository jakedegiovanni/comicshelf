package main

import (
	"html/template"
	"log/slog"
	"net/http"

	"github.com/jakedegiovanni/comicshelf/static"
)

const SeriesEndpoint = "/marvel-unlimited/series"

type Series struct {
	tmpl   *template.Template
	client *MarvelClient
	db     *Db
	logger *slog.Logger
}

func NewSeries(tmpl *template.Template, client *MarvelClient, db *Db, logger *slog.Logger) *Series {
	return &Series{
		tmpl:   tmpl,
		client: client,
		db:     db,
		logger: logger,
	}
}

func (s *Series) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.logger.Debug(r.URL.String())
	s.logger.Debug(r.URL.Query().Get("series"))

	resp, err := s.client.GetComicsWithinSeries(r.URL.Query().Get("series"))
	if err != nil {
		s.logger.Error(err.Error())
		return
	}

	content := static.Content{
		Date:         r.URL.Query().Get("date"),
		PageEndpoint: SeriesEndpoint,
		Resp:         resp,
	}

	err = s.tmpl.ExecuteTemplate(w, "index.html", content)
	if err != nil {
		s.logger.Error(err.Error())
		return
	}
}
