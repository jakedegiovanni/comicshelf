package main

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

//go:embed index.html
var index string

//go:embed static
var static embed.FS

func main() {
	db, err := NewDb("db.json")
	if err != nil {
		log.Fatalln(err)
	}

	defer db.Shutdown()

	client := NewMarvelClient()

	tmpl := template.Must(
		template.
			New("tmpl").
			Funcs(template.FuncMap{
				"following": db.Following,
			}).
			Parse(index),
	)

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodPost {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		if r.Method == http.MethodPost {
			_ = r.ParseForm()
			if r.PostForm.Has("follow") {
				db.Follow(r.PostFormValue("follow"))
			} else if r.PostForm.Has("unfollow") {
				db.Unfollow(r.PostFormValue("unfollow"))
			} else {
				log.Println("unknown postform values")
			}

			http.Redirect(w, r, fmt.Sprintf("/?date=%s", r.PostFormValue("date")), http.StatusFound)
			return
		}

		if !r.URL.Query().Has("date") {
			http.Redirect(w, r, fmt.Sprintf("/?date=%s", time.Now().Format("2006-01-02")), http.StatusFound)
			return
		}

		t, err := time.Parse("2006-01-02", r.URL.Query().Get("date"))
		if err != nil {
			log.Fatalln("parse", err)
		}

		resp, err := client.GetWeeklyComics(t)
		if err != nil {
			log.Fatalln(fmt.Errorf("getting series collection: %w", err))
		}

		input := make(map[string]interface{})
		input["resp"] = resp
		input["date"] = r.URL.Query().Get("date")

		err = tmpl.Execute(w, input)
		if err != nil {
			log.Fatalln("exec", err)
		}
	})

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
