package main

import (
	"embed"
	"errors"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"strings"
	"syscall"
)

//go:embed static
var static embed.FS

type Content struct {
	Date         string
	PageEndpoint string
	Resp         interface{}
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}))

	db, err := NewDb("db.json", logger)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer HandlePanic(logger)
	defer db.Shutdown()

	client := NewMarvelClient(logger)

	tmpl := template.Must(
		template.
			New("marvel-unlimited").
			Funcs(template.FuncMap{
				"contains": func(s1, s2 string) bool {
					return strings.Contains(strings.ToLower(s1), strings.ToLower(s2))
				},
				"equals":    strings.EqualFold,
				"following": db.Following,
				"content": func(vals ...interface{}) (map[interface{}]interface{}, error) {
					if len(vals)%2 != 0 {
						return nil, errors.New("invalid dict call")
					}

					dict := make(map[interface{}]interface{})
					for i := 0; i < len(vals); i += 2 {
						dict[vals[i]] = vals[i+1]
					}
					return dict, nil
				},
			}).
			ParseFS(static, "**/index.html", "**/marvel-unlimited.html", "**/comic-card.html"),
	)

	comics := NewComics(tmpl, client, db, logger)
	series := NewSeries(tmpl, client, db, logger)

	mux := http.NewServeMux()

	chain := MiddlewareChain(
		RecoverHandler(logger),
		AllowedMethods(http.MethodGet, http.MethodPost),
	)

	mux.HandleFunc(ComicsEndpoint, chain(comics.ServeHTTP))
	mux.HandleFunc(SeriesEndpoint, chain(series.ServeHTTP))

	f, err := fs.Sub(static, "static")
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(f))))

	srv := &http.Server{
		Handler: mux,
		Addr:    "127.0.0.1:8080",
	}

	go func() {
		err = srv.ListenAndServe()
		if err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}
	}()

	logger.Info("server ready to accept connections")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
}

func HandlePanic(logger *slog.Logger) {
	if r := recover(); r != nil {
		logger.Error("recovered", slog.Any("r", r))
		debug.PrintStack()
	}
}
