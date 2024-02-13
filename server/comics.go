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
	router.Get("/marvel-unlimited", s.handleMarvelUnlimitedComics)
}

func (s *Server) handleMarvelUnlimitedComics(w http.ResponseWriter, r *http.Request) {
	if !r.URL.Query().Has("date") {
		query := r.URL.Query()
		query.Set("date", time.Now().Format("2006-01-02"))
		r.URL.RawQuery = query.Encode()

		http.Redirect(w, r, r.URL.String(), http.StatusFound)
		return
	}

	t, err := time.Parse("2006-01-02", r.URL.Query().Get("date"))
	if err != nil {
		s.logger.Error(err.Error())
		return
	}

	resp, err := s.comics.GetWeeklyComics(r.Context(), t)
	if err != nil {
		s.logger.Error("getting series collection", slog.String("err", err.Error()))
		return
	}

	content := templates.View[comicshelf.Page[comicshelf.Comic]]{
		Date: r.URL.Query().Get("date"),
		Resp: resp,
	}

	err = s.tmpl.ExecuteTemplate(w, "index.html", content)
	if err != nil {
		s.logger.Error(err.Error())
	}
}
