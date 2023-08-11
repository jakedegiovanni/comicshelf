package main

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"
)

const ComicsEndpoint = "/marvel-unlimited/comics/"

type Comics struct {
	tmpl   *template.Template
	client *MarvelClient
	db     *Db
	logger *slog.Logger
}

func NewComics(tmpl *template.Template, client *MarvelClient, db *Db, logger *slog.Logger) *Comics {
	return &Comics{
		tmpl:   tmpl,
		client: client,
		db:     db,
		logger: logger,
	}
}

func (c *Comics) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.URL.Path, "track") {
		_ = r.ParseForm()
		id := r.PostFormValue("id")
		name := r.PostFormValue("name")
		r.Header.Set("Content-Type", "text/html; charset=utf-8")
		if c.db.Following(id) {
			c.db.Unfollow(id)
			_, err := w.Write([]byte(`<svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg" style="width: 32px;height: 32px;"><g id="SVGRepo_bgCarrier" stroke-width="0"></g><g id="SVGRepo_tracerCarrier" stroke-linecap="round" stroke-linejoin="round"></g><g id="SVGRepo_iconCarrier"> <path fill-rule="evenodd" clip-rule="evenodd" d="M4 4C4 2.34315 5.34315 1 7 1H17C18.6569 1 20 2.34315 20 4V20.9425C20 22.6114 18.0766 23.5462 16.7644 22.5152L12 18.7717L7.23564 22.5152C5.92338 23.5462 4 22.6114 4 20.9425V4ZM7 3C6.44772 3 6 3.44772 6 4V20.9425L12 16.2283L18 20.9425V4C18 3.44772 17.5523 3 17 3H7Z" fill="#0F0F0F"></path> </g></svg>`))
			if err != nil {
				c.logger.Warn("error writing unfollow", slog.String("err", err.Error()))
			}
			return
		} else {
			c.db.Follow(id, name)
			_, err := w.Write([]byte(`<svg viewBox="-4 0 30 30" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" xmlns:sketch="http://www.bohemiancoding.com/sketch/ns" fill="#000000"><g id="SVGRepo_bgCarrier" stroke-width="0"></g><g id="SVGRepo_tracerCarrier" stroke-linecap="round" stroke-linejoin="round"></g><g id="SVGRepo_iconCarrier"> <title>bookmark</title> <desc>Created with Sketch Beta.</desc> <defs> </defs> <g id="Page-1" stroke="none" stroke-width="1" fill="none" fill-rule="evenodd" sketch:type="MSPage"> <g id="Icon-Set-Filled" sketch:type="MSLayerGroup" transform="translate(-419.000000, -153.000000)" fill="#000000"> <path d="M437,153 L423,153 C420.791,153 419,154.791 419,157 L419,179 C419,181.209 420.791,183 423,183 L430,176 L437,183 C439.209,183 441,181.209 441,179 L441,157 C441,154.791 439.209,153 437,153" id="bookmark" sketch:type="MSShapeGroup"> </path> </g> </g> </g></svg>`))
			if err != nil {
				c.logger.Warn("error writing follow", slog.String("err", err.Error()))
			}
			return
		}
	}

	if !r.URL.Query().Has("date") {
		http.Redirect(w, r, fmt.Sprintf("/marvel-unlimited/comics?date=%s", time.Now().Format("2006-01-02")), http.StatusFound)
		return
	}

	t, err := time.Parse("2006-01-02", r.URL.Query().Get("date"))
	if err != nil {
		c.logger.Error("parse", slog.String("err", err.Error()))
		os.Exit(1) // todo - shouldn't be doing this
	}

	resp, err := c.client.GetWeeklyComics(t)
	if err != nil {
		c.logger.Error("getting series collection", slog.String("err", err.Error()))
		os.Exit(1) // todo shouldn't bee doing this
	}

	content := Content{
		Date:         r.URL.Query().Get("date"),
		PageEndpoint: ComicsEndpoint,
		Resp:         resp,
	}

	err = c.tmpl.ExecuteTemplate(w, "index.html", content)
	if err != nil {
		c.logger.Error("exec", slog.String("err", err.Error()))
		os.Exit(1) // todo - shouldn't be doing this
	}
}
