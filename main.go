package main

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"time"
)

//go:embed index.html
var index string

//go:embed static
var static embed.FS

func main() {
	client := NewMarvelClient()

	tmpl := template.Must(
		template.
			New("tmpl").
			Funcs(template.FuncMap{
				"following": following,
			}).
			Parse(index),
	)

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodPost {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
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

		if r.Method == http.MethodPost {
			r.ParseForm()
			if r.PostForm.Has("follow") {
				follow(r.PostFormValue("follow"))
			} else if r.PostForm.Has("unfollow") {
				unfollow(r.PostFormValue("unfollow"))
			} else {
				log.Println("unknown postform values")
			}
		}

		input := make(map[string]interface{})
		input["resp"] = resp
		input["week"] = r.URL.Query().Get("date")

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

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}
