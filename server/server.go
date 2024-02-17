package server

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jakedegiovanni/comicshelf"
	"github.com/jakedegiovanni/comicshelf/server/static"
	"github.com/jakedegiovanni/comicshelf/server/templates"
	"golang.org/x/sync/errgroup"
)

const justTheDateFormat = "2006-01-02"

type Server struct {
	cfg    *Config
	srv    *http.Server
	logger *slog.Logger
	tmpl   *template.Template
	comics comicshelf.ComicService
	series comicshelf.SeriesService
	user   comicshelf.UserService
}

func New(
	config *Config, logger *slog.Logger,
	comics comicshelf.ComicService,
	series comicshelf.SeriesService,
	user comicshelf.UserService,
) *Server {
	router := chi.NewRouter()

	srv := &http.Server{
		Handler: router,
		Addr:    config.Address,
	}

	tmpl := template.Must(
		template.
			New("comicshelf").
			Funcs(template.FuncMap{
				"equals": strings.EqualFold,
			}).
			ParseFS(templates.Files, "*.html"),
	)

	s := &Server{
		cfg:    config,
		srv:    srv,
		logger: logger,
		tmpl:   tmpl,
		comics: comics,
		series: series,
		user:   user,
	}

	router.Use(serverLogger(logger))
	router.Use(middleware.Recoverer)

	router.Mount("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(static.Files))))

	router.Route("/comics", func(r chi.Router) {
		s.registerComicRoutes(r)
	})

	router.Route("/series", func(r chi.Router) {
		s.registerSeriesRoutes(r)
	})

	router.Route("/api", func(r chi.Router) {
		s.registerUserRoutes(r)
	})

	return s
}

func (s *Server) Run(ctx context.Context) error {
	defer s.handlePanic()

	g := new(errgroup.Group)

	g.Go(func() error {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)

		select {
		case <-ctx.Done():
			s.logger.Info("programmatic shutdown")
		case <-c:
			s.logger.Info("signal shutdown")
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		return s.srv.Shutdown(ctx)
	})

	g.Go(func() error {
		err := s.srv.ListenAndServe()
		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return nil
			}

			s.logger.Error(err.Error())
			return fmt.Errorf("error starting server: %w", err)
		}

		return nil
	})

	s.logger.Info("server ready to accept connections", slog.String("addr", s.cfg.Address))
	return g.Wait()
}

func (s *Server) handlePanic() {
	if r := recover(); r != nil {
		s.logger.Error("recovered", slog.Any("r", r))
		debug.PrintStack()
	}
}
