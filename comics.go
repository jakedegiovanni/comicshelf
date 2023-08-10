package main

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"
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
	if r.Method == http.MethodPost {
		_ = r.ParseForm()
		if r.PostForm.Has("follow") {
			c.db.Follow(r.PostFormValue("follow"), r.PostFormValue("name"))
		} else if r.PostForm.Has("unfollow") {
			c.db.Unfollow(r.PostFormValue("unfollow"))
		} else {
			c.logger.Warn("unknown postform values")
		}

		http.Redirect(w, r, fmt.Sprintf("/marvel-unlimited/comics?date=%s", r.PostFormValue("date")), http.StatusFound)
		return
	}

	if !r.URL.Query().Has("date") {
		http.Redirect(w, r, fmt.Sprintf("/marvel-unlimited/comics?date=%s", time.Now().Format("2006-01-02")), http.StatusFound)
		return
	}

	t, err := time.Parse("2006-01-02", r.URL.Query().Get("date"))
	if err != nil {
		c.logger.Error("parse", slog.String("err", err.Error()))
		os.Exit(1) // todo - shouldn't be doing this
	}

	resp, err := c.client.GetWeeklyComics(t)
	if err != nil {
		c.logger.Error("getting series collection", slog.String("err", err.Error()))
		os.Exit(1) // todo shouldn't bee doing this
	}

	content := Content{
		Date:         r.URL.Query().Get("date"),
		PageEndpoint: ComicsEndpoint,
		Resp:         resp,
	}

	err = c.tmpl.ExecuteTemplate(w, "index.html", content)
	if err != nil {
		c.logger.Error("exec", slog.String("err", err.Error()))
		os.Exit(1) // todo - shouldn't be doing this
	}
}
