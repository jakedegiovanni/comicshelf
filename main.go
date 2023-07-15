package main

import (
	"crypto/md5"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"time"
)

//go:embed index.html
var index string

//go:embed static
var static embed.FS

const (
	pub  = "e44b60fc0a6fb0fd3bcc483c8be2df3f"
	priv = "802447629feca7506db62598688e9c597815ece0"
)

type apiKeyMiddleWare struct {
	next http.RoundTripper
}

func (a *apiKeyMiddleWare) RoundTrip(req *http.Request) (*http.Response, error) {
	ts := fmt.Sprintf("%d", time.Now().UTC().Unix())
	hash := md5.Sum([]byte(ts + priv + pub))
	query := req.URL.Query()
	query.Add("ts", ts)
	query.Add("hash", fmt.Sprintf("%x", hash))
	query.Add("apikey", pub)
	req.URL.RawQuery = query.Encode()
	return a.next.RoundTrip(req)
}

type addBase struct {
	next http.RoundTripper
}

func (a *addBase) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Host = "gateway.marvel.com"
	req.URL.Host = req.Host
	req.URL.Scheme = "https"
	req.URL.Path = fmt.Sprintf("/v1/public%s", req.URL.Path)
	return a.next.RoundTrip(req)
}

func weekRange(t time.Time) (time.Time, time.Time) {
	for t.Weekday() != time.Sunday {
		t = t.AddDate(0, 0, -1)
	}

	t = t.AddDate(0, 0, -7)

	return t, t.AddDate(0, 0, 6)
}

func main() {
	client := &http.Client{
		Timeout:   20 * time.Second,
		Transport: &addBase{next: &apiKeyMiddleWare{next: http.DefaultTransport}},
	}

	tmpl := template.New("tmpl")

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Has("marvel") {
			first, last := weekRange(time.Now().AddDate(0, -3, 0))
			layout := "2006-01-02"
			resp, err := client.Get(fmt.Sprintf("/comics?format=comic&formatType=comic&noVariants=true&dateRange=%s,%s&hasDigitalIssue=true&orderBy=issueNumber&limit=100", first.Format(layout), last.Format(layout)))
			if err != nil {
				log.Fatalln(fmt.Errorf("getting series collection: %w", err))
			}

			log.Printf("%+v\n", resp.Header)

			defer resp.Body.Close()

			f, err := os.OpenFile("results.json", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
			if err != nil {
				log.Fatalln(fmt.Errorf("opening file: %w", err))
			}
			defer f.Close()

			_, err = io.Copy(f, resp.Body)
			if err != nil {
				log.Fatalln(fmt.Errorf("writing: %w", err))
			}

			log.Println("ok")
		}
		f, err := os.OpenFile("results.json", os.O_CREATE|os.O_RDONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalln(err)
		}
		defer f.Close()

		var resp map[string]interface{}
		err = json.NewDecoder(f).Decode(&resp)
		if err != nil {
			log.Fatalln("decode", err)
		}

		t, err := tmpl.Clone()
		if err != nil {
			log.Fatalln("clone", err)
		}

		t, err = t.Parse(index)
		if err != nil {
			log.Fatalln("parse", err)
		}

		err = t.Execute(w, resp)
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
