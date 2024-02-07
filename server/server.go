package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type ServerConfig struct {
	Address string
}

func NewServer(config ServerConfig) (*http.Server, *chi.Mux) {
	router := chi.NewRouter()

	router.Use(middleware.Recoverer)

	return &http.Server{
		Handler: router,
		Addr:    config.Address,
	}, router
}
