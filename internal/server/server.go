package server

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
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
	"golang.org/x/sync/errgroup"
)

const justTheDateFormat = "2006-01-02"

//go:embed static
var static embed.FS

//go:embed templates
var templates embed.FS

type View[T any] struct {
	Date  string
	Title string
	Resp  T
}

type Server struct {
	cfg        *Config
	srv        *http.Server
	comicTmpl  *template.Template
	seriesTmpl *template.Template
	comics     comicshelf.ComicService
	series     comicshelf.SeriesService
	user       comicshelf.UserService
}

func New(
	config *Config,
	comics comicshelf.ComicService,
	series comicshelf.SeriesService,
	user comicshelf.UserService,
) (*Server, error) {
	router := chi.NewRouter()

	srv := &http.Server{
		Handler: router,
		Addr:    config.Address,
	}

	tmplFuncs := template.FuncMap{
		"equals": strings.EqualFold,
		"following": func(userId, seriesId int) bool {
			// todo this shouldn't stay here when an actual db connection, don't want to be calling sequentially during template render
			f, err := user.Following(context.TODO(), userId, seriesId)
			if err != nil {
				slog.Error(err.Error())
				return false
			}
			slog.Debug(fmt.Sprintf("%+v", f))
			return f
		},
		"justTheDate": func(t time.Time) string {
			return t.Format(justTheDateFormat)
		},
	}

	templates, err := fs.Sub(templates, "templates")
	if err != nil {
		return nil, err
	}

	comicTmpl := template.Must(
		template.
			New("comicTmpl").
			Funcs(tmplFuncs).
			ParseFS(templates, "*.html", "comics/*.html"),
	)

	seriesTmpl := template.Must(
		template.
			New("seriesTmpl").
			Funcs(tmplFuncs).
			ParseFS(templates, "*.html", "series/*.html"),
	)

	s := &Server{
		cfg:        config,
		srv:        srv,
		comicTmpl:  comicTmpl,
		seriesTmpl: seriesTmpl,
		comics:     comics,
		series:     series,
		user:       user,
	}

	router.Use(serverLogger())
	router.Use(middleware.Recoverer)

	router.Group(func(r chi.Router) {
		r.Mount("/static/", http.FileServer(http.FS(static)))
	})

	router.Group(func(r chi.Router) {
		r.Use(middleware.StripSlashes)

		r.Route("/comics", func(r chi.Router) {
			s.registerComicRoutes(r)
		})

		r.Route("/series", func(r chi.Router) {
			s.registerSeriesRoutes(r)
		})

		r.Route("/api", func(r chi.Router) {
			s.registerUserRoutes(r)
		})
	})

	return s, nil
}

func (s *Server) Run(ctx context.Context) error {
	defer s.handlePanic()

	g := new(errgroup.Group)

	g.Go(func() error {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)

		select {
		case <-ctx.Done():
			slog.Info("programmatic shutdown")
		case <-c:
			slog.Info("signal shutdown")
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

			slog.Error(err.Error())
			return fmt.Errorf("error starting server: %w", err)
		}

		return nil
	})

	slog.Info("server ready to accept connections", slog.String("addr", s.cfg.Address))
	return g.Wait()
}

func (s *Server) handlePanic() {
	if r := recover(); r != nil {
		slog.Error("recovered", slog.Any("r", r))
		debug.PrintStack()
	}
}
