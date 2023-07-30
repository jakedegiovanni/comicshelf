package main

import (
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"time"
)

const ComicsEndpoint = "/marvel-unlimited/comics"

type Comics struct {
	tmpl   *template.Template
	client *MarvelClient
	db     *Db
}

func NewComics(sys fs.FS, client *MarvelClient, db *Db) *Comics {
	tmpl := template.Must(
		template.
			New("comics").
			Funcs(template.FuncMap{
				"following": db.Following,
			}).
			ParseFS(sys, "**/index.html", "**/comics/comics.html"),
	)

	return &Comics{
		tmpl:   tmpl,
		client: client,
		db:     db,
	}
}

func (c *Comics) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	if r.Method == http.MethodPost {
		_ = r.ParseForm()
		if r.PostForm.Has("follow") {
			c.db.Follow(r.PostFormValue("follow"))
		} else if r.PostForm.Has("unfollow") {
			c.db.Unfollow(r.PostFormValue("unfollow"))
		} else {
			log.Println("unknown postform values")
		}

		http.Redirect(w, r, fmt.Sprintf("/marvel-unlimited/comics?date=%s", r.PostFormValue("date")), http.StatusFound)
		return
	}

	if !r.URL.Query().Has("date") {
		http.Redirect(w, r, fmt.Sprintf("/marvel-unlimited/comics?date=%s", time.Now().Format("2006-01-02")), http.StatusFound)
		return
	}

	t, err := time.Parse("2006-01-02", r.URL.Query().Get("date"))
	if err != nil {
		log.Fatalln("parse", err)
	}

	resp, err := c.client.GetWeeklyComics(t)
	if err != nil {
		log.Fatalln(fmt.Errorf("getting series collection: %w", err))
	}

	content := Content{
		Date:         r.URL.Query().Get("date"),
		PageEndpoint: ComicsEndpoint,
		Resp:         resp,
	}

	err = c.tmpl.ExecuteTemplate(w, "index.html", content)
	if err != nil {
		log.Fatalln("exec", err)
	}
}
