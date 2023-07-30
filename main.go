package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
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

	comics := NewComics(static, client, db)

	mux := http.NewServeMux()

	mux.HandleFunc(ComicsEndpoint, RecoverHandler(comics.ServeHTTP))

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
