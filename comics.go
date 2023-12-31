package main

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"time"

	"github.com/jakedegiovanni/comicshelf/static"
)

const ComicsEndpoint = "/marvel-unlimited/comics"

type Comics struct {
	tmpl   *template.Template
	client *MarvelClient
	db     *Db
	logger *slog.Logger
}

func NewComics(tmpl *template.Template, client *MarvelClient, db *Db, logger *slog.Logger) *Comics {
	return &Comics{
		tmpl:   tmpl,
		client: client,
		db:     db,
		logger: logger,
	}
}

func (c *Comics) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !r.URL.Query().Has("date") {
		http.Redirect(w, r, fmt.Sprintf("/marvel-unlimited/comics?date=%s", time.Now().Format("2006-01-02")), http.StatusFound)
		return
	}

	t, err := time.Parse("2006-01-02", r.URL.Query().Get("date"))
	if err != nil {
		c.logger.Error(err.Error())
		return
	}

	resp, err := c.client.GetWeeklyComics(t)
	if err != nil {
		c.logger.Error("getting series collection", slog.String("err", err.Error()))
		return
	}

	content := static.Content{
		Date:         r.URL.Query().Get("date"),
		PageEndpoint: ComicsEndpoint,
		Resp:         resp,
	}

	err = c.tmpl.ExecuteTemplate(w, "index.html", content)
	if err != nil {
		c.logger.Error(err.Error())
	}
}
