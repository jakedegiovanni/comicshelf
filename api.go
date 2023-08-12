package main

import (
	"html/template"
	"log/slog"
	"net/http"
)

const TrackEndpoint = "/api/track"

type Api struct {
	logger *slog.Logger
	db     *Db
	tmpl   *template.Template
}

func NewApi(logger *slog.Logger, db *Db, tmpl *template.Template) *Api {
	return &Api{
		logger: logger,
		db:     db,
		tmpl:   tmpl,
	}
}

func (a *Api) Track(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	key := r.PostFormValue("key")
	name := r.PostFormValue("name")
	if a.db.Following(key) {
		a.db.Unfollow(key)
		err := a.tmpl.ExecuteTemplate(w, "follow", nil)
		if err != nil {
			a.logger.Warn("error writing unfollow", slog.String("err", err.Error()))
		}
		return
	}

	a.db.Follow(key, name)
	err := a.tmpl.ExecuteTemplate(w, "unfollow", nil)
	if err != nil {
		a.logger.Warn("error writing follow", slog.String("err", err.Error()))
	}
}
