package server

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jakedegiovanni/comicshelf"
	"github.com/jakedegiovanni/comicshelf/server/templates"
)

func (s *Server) registerSeriesRoutes(router chi.Router) {
	router.Get("/{seriesId}", s.handleSeries)
}

func (s *Server) handleSeries(w http.ResponseWriter, r *http.Request) {
	slog.Debug(r.URL.String())

	seriesId := chi.URLParam(r, "seriesId")
	id, err := strconv.Atoi(seriesId)
	if err != nil {
		http.Error(w, "series query is not a number", http.StatusUnprocessableEntity)
		return
	}

	resp, err := s.series.GetComicsWithinSeries(r.Context(), id)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	content := templates.View[[]comicshelf.Comic]{
		Date:  r.URL.Query().Get("date"),
		Title: "Series Issues",
		Resp:  resp,
	}

	err = s.seriesTmpl.ExecuteTemplate(w, "index.html", content)
	if err != nil {
		slog.Error(err.Error())
		return
	}
}
