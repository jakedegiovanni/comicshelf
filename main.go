package main

import (
	"embed"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"strings"
)

//go:embed static
var static embed.FS

type Content struct {
	Date         string
	PageEndpoint string
	Resp         interface{}
}

func main() {
	db, err := NewDb("db.json")
	if err != nil {
		log.Fatalln(err)
	}

	defer HandlePanic()
	defer db.Shutdown()

	client := NewMarvelClient()

	tmpl := template.Must(
		template.
			New("marvel-unlimited").
			Funcs(template.FuncMap{
				"contains": func(s1, s2 string) bool {
					return strings.Contains(strings.ToLower(s1), strings.ToLower(s2))
				},
				"equals": func(s1, s2 string) bool {
					return strings.ToLower(s1) == strings.ToLower(s2)
				},
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

	comics := NewComics(tmpl, client, db)
	series := NewSeries(tmpl, client, db)

	mux := http.NewServeMux()

	chain := MiddlewareChain(
		RecoverHandler(),
		AllowedMethods(http.MethodGet, http.MethodPost),
	)

	mux.HandleFunc(ComicsEndpoint, chain(comics.ServeHTTP))
	mux.HandleFunc(SeriesEndpoint, chain(series.ServeHTTP))

	f, err := fs.Sub(static, "static")
	if err != nil {
		log.Fatalln(err)
	}
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(f))))

	srv := &http.Server{
		Handler: mux,
		Addr:    "127.0.0.1:8080",
	}

	go func() {
		err = srv.ListenAndServe()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	fmt.Println("server ready to accept connections")
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill)

	<-c
}

func HandlePanic() {
	if r := recover(); r != nil {
		log.Println("recovered", r)
		debug.PrintStack()
	}
}
