package server

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jakedegiovanni/comicshelf"
	"github.com/jakedegiovanni/comicshelf/server/templates"
)

func (s *Server) registerComicRoutes(router chi.Router) {
	router.Use(queryDate())
	router.Get("/marvel-unlimited", s.handleMarvelUnlimitedComics)
}

func (s *Server) handleMarvelUnlimitedComics(w http.ResponseWriter, r *http.Request) {
	t, err := time.Parse(justTheDateFormat, r.URL.Query().Get("date"))
	if err != nil {
		slog.Error(err.Error())
		return
	}

	comics, err := s.comics.GetWeeklyComics(r.Context(), t)
	if err != nil {
		slog.Error("getting series collection", slog.String("err", err.Error()))
		return
	}

	content := templates.View[comicshelf.Page[comicshelf.Comic]]{
		Date: r.URL.Query().Get("date"),
		Resp: comics,
	}

	err = s.tmpl.ExecuteTemplate(w, "index.html", content)
	if err != nil {
		slog.Error(err.Error())
	}
}
