package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (s *Server) registerUserRoutes(router chi.Router) {
	router.Post("/api/follow", s.registerFollow)
	router.Post("/api/unfollow", s.registerUnfollow)
}

func (s *Server) registerFollow(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	_ = r.PostFormValue("key")
	_ = r.PostFormValue("name")
	// if a.db.Following(key) {
	// 	a.db.Unfollow(key)
	// 	err := a.tmpl.ExecuteTemplate(w, "follow", nil)
	// 	if err != nil {
	// 		a.logger.Warn("error writing unfollow", slog.String("err", err.Error()))
	// 	}
	// 	return
	// }

	// a.db.Follow(key, name)
	// err := a.tmpl.ExecuteTemplate(w, "unfollow", nil)
	// if err != nil {
	// 	a.logger.Warn("error writing follow", slog.String("err", err.Error()))
	// }
}

func (s *Server) registerUnfollow(w http.ResponseWriter, r *http.Request) {}
