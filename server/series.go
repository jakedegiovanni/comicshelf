package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jakedegiovanni/comicshelf"
	"github.com/jakedegiovanni/comicshelf/server/templates"
)

func (s *Server) registerSeriesRoutes(router chi.Router) {
	router.Get("/marvel-unlimited", s.handleMarvelUnlimitedSeries)
}

func (s *Server) handleMarvelUnlimitedSeries(w http.ResponseWriter, r *http.Request) {
	s.logger.Debug(r.URL.String())
	s.logger.Debug(r.URL.Query().Get("series"))

	// resp, err := s.series.GetComicsWithinSeries(r.Context(), r.URL.Query().Get("series"))
	resp, err := s.series.GetComicsWithinSeries(r.Context(), 0)
	if err != nil {
		s.logger.Error(err.Error())
		return
	}

	content := templates.View[[]comicshelf.Comic]{
		Date: r.URL.Query().Get("date"),
		Resp: resp,
	}

	err = s.tmpl.ExecuteTemplate(w, "index.html", content)
	if err != nil {
		s.logger.Error(err.Error())
		return
	}
}
