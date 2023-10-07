package main

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
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jakedegiovanni/comicshelf/static"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
)

func Server(cfg *AppConfig, v *viper.Viper) *cobra.Command {
	cmd := &cobra.Command{
		Use: "server",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := GetConfigFromCtx(cmd.Context())
			if err != nil {
				return err
			}

			logger, err := cfg.Logging.Logger()
			if err != nil {
				return err
			}

			db, err := NewDb("db.json", logger)
			if err != nil {
				logger.Error(err.Error())
				return err
			}

			defer db.Shutdown()
			defer HandlePanic(logger)

			client := NewMarvelClient(logger)

			tmpl := template.Must(
				template.
					New("marvel-unlimited").
					Funcs(template.FuncMap{
						"contains": func(s1, s2 string) bool {
							return strings.Contains(strings.ToLower(s1), strings.ToLower(s2))
						},
						"equals":              strings.EqualFold,
						"following":           db.Following,
						"marvelUnlimitedDate": DateResponseToMarvelUnlimitedDate,
					}).
					ParseFS(static.Files, "index.html", "marvel-unlimited.html", "comic-card.html", "follow.html", "unfollow.html"),
			)

			comics := NewComics(tmpl, client, db, logger)
			series := NewSeries(tmpl, client, db, logger)
			api := NewApi(logger, db, tmpl)

			router := chi.NewRouter()

			router.Use(ServerLogger(logger))
			router.Use(middleware.Recoverer)

			router.Get(ComicsEndpoint, comics.ServeHTTP)
			router.Get(SeriesEndpoint, series.ServeHTTP)
			router.Post(TrackEndpoint, api.Track)

			router.Mount("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(static.Files))))

			srv := &http.Server{
				Handler: router,
				Addr:    cfg.Server.Address,
			}

			g := new(errgroup.Group)

			g.Go(func() error {
				c := make(chan os.Signal, 1)
				signal.Notify(c, os.Interrupt, syscall.SIGTERM)

				select {
				case <-cmd.Context().Done():
					logger.Info("context shutdown")
				case <-c:
					logger.Info("programmatic shutdown")
				}

				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				return srv.Shutdown(ctx)
			})

			g.Go(func() error {
				err = srv.ListenAndServe()
				if err != nil {
					if errors.Is(err, http.ErrServerClosed) {
						return nil
					}

					logger.Error(err.Error())
					return fmt.Errorf("error starting server: %w", err)
				}
				return nil
			})

			logger.Info("server ready to accept connections", slog.String("addr", cfg.Server.Address))

			return g.Wait()
		},
	}

	cmd.PersistentFlags().StringVarP(&cfg.Server.Address, "address", "a", cfg.Server.Address, "")
	v.BindPFlag("server.address", cmd.PersistentFlags().Lookup("address"))

	return cmd
}

func HandlePanic(logger *slog.Logger) {
	if r := recover(); r != nil {
		logger.Error("recovered", slog.Any("r", r))
		debug.PrintStack()
	}
}
