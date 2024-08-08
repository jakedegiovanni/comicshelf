package server

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (s *Server) registerUserRoutes(router chi.Router) {
	router.Post("/follow", s.registerFollow)
	router.Post("/unfollow", s.registerUnfollow)
}

func (s *Server) registerFollow(w http.ResponseWriter, r *http.Request) {
	id, err := s.extractSeriesIdFromForm(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not extract series id: %s", err.Error()), http.StatusBadRequest)
		return
	}
	slog.Debug(fmt.Sprintf("%d", id))

	err = s.user.Follow(r.Context(), 0, id) // using default user id until auth actually implemented
	if err != nil {
		http.Error(w, fmt.Sprintf("could not follow series with id: %d", id), http.StatusInternalServerError)
		return
	}

	err = s.comicTmpl.ExecuteTemplate(w, "unfollow", nil)
	if err != nil {
		slog.Warn("error writing unfollow", slog.String("err", err.Error()))
	}
}

func (s *Server) registerUnfollow(w http.ResponseWriter, r *http.Request) {
	id, err := s.extractSeriesIdFromForm(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not extract series id: %s", err.Error()), http.StatusBadRequest)
		return
	}

	err = s.user.Unfollow(r.Context(), 0, id) // using default user id until auth actually implemented
	if err != nil {
		http.Error(w, fmt.Sprintf("could not unfollow series with id: %d", id), http.StatusInternalServerError)
		return
	}

	err = s.comicTmpl.ExecuteTemplate(w, "follow", nil)
	if err != nil {
		slog.Warn("error writing follow", slog.String("err", err.Error()))
	}
}

func (s *Server) extractSeriesIdFromForm(r *http.Request) (int, error) {
	err := r.ParseForm()
	if err != nil {
		return 0, errors.New("could not read form")
	}

	seriesId := r.PostFormValue("series")
	if seriesId == "" {
		return 0, errors.New("series key not present")
	}

	id, err := strconv.Atoi(seriesId)
	if err != nil {
		return 0, fmt.Errorf("series id is not a valid number: %s", seriesId)
	}

	return id, nil
}
