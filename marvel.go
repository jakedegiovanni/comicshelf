package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

const (
	monthOffset = -3
	layout      = "2006-01-02"
)

type apiKeyMiddleWare struct {
	next http.RoundTripper
	pub  io.ReadSeeker
	priv io.ReadSeeker
}

func (a *apiKeyMiddleWare) RoundTrip(req *http.Request) (*http.Response, error) {
	_, err := a.pub.Seek(0, io.SeekStart)
	if err != nil {
		log.Fatalln(err)
	}

	pub, err := io.ReadAll(a.pub)
	if err != nil {
		log.Fatalln("pub read", err)
	}

	_, err = a.priv.Seek(0, io.SeekStart)
	if err != nil {
		log.Fatalln(err)
	}

	priv, err = io.ReadAll(a.priv)
	if err != nil {
		log.Fatalln("priv read", err)
	}

	ts := fmt.Sprintf("%d", time.Now().UTC().Unix())
	hash := md5.Sum([]byte(ts + string(priv) + string(pub)))
	query := req.URL.Query()
	query.Add("ts", ts)
	query.Add("hash", fmt.Sprintf("%x", hash))
	query.Add("apikey", string(pub))
	req.URL.RawQuery = query.Encode()
	return a.next.RoundTrip(req)
}

type DataWrapper struct {
	Code            interface{}   `json:"code"`
	Status          string        `json:"status"`
	Copyright       string        `json:"copyright"`
	AttributionText string        `json:"attributionText"`
	AttributionHTML string        `json:"attributionHTML"`
	Etag            string        `json:"etag"`
	Data            DataContainer `json:"data"`
}

type DataContainer struct {
	Offset  int                      `json:"offset"`
	Limit   int                      `json:"limit"`
	Total   int                      `json:"total"`
	Count   int                      `json:"count"`
	Results []map[string]interface{} `json:"results"`
}

type MarvelClient struct {
	client    *http.Client
	etagCache map[string]*DataWrapper
	mu        *sync.Mutex
}

func NewMarvelClient() *MarvelClient {
	return &MarvelClient{
		client: &http.Client{
			Timeout: 20 * time.Second,
			Transport: &addBase{
				next: &apiKeyMiddleWare{
					next: http.DefaultTransport,
					pub:  Pub,
					priv: Priv,
				},
			},
		},
		etagCache: make(map[string]*DataWrapper),
		mu:        &sync.Mutex{},
	}
}

func (m *MarvelClient) GetWeeklyComics(t time.Time) (*DataWrapper, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if t.Weekday() == time.Sunday {
		// sunday is when the cut off is however new comics show up on the monday, the offset set here makes expectations match reality
		t = t.AddDate(0, 0, -1)
	}

	first, last := m.weekRange(t.AddDate(0, monthOffset, 0))
	endpoint := fmt.Sprintf("/comics?format=comic&formatType=comic&noVariants=true&dateRange=%s,%s&hasDigitalIssue=true&orderBy=issueNumber&limit=100", first.Format(layout), last.Format(layout))

	if resp, ok := m.etagCache[endpoint]; ok {
		req, err := http.NewRequest(http.MethodGet, endpoint, nil)
		if err != nil {
			return nil, err
		}

		req.Header.Set("If-None-Match", resp.Etag)

		r, err := m.client.Do(req)
		if err != nil {
			return nil, err
		}
		defer r.Body.Close()

		if r.StatusCode == http.StatusNotModified {
			log.Println("not modified, using cached response")
			return resp, nil
		}
	}

	log.Println("item modified or not present in cache")

	resp, err := m.client.Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var d DataWrapper
	err = json.NewDecoder(resp.Body).Decode(&d)
	if err != nil {
		return nil, err
	}

	m.etagCache[endpoint] = &d
	return &d, nil
}

func (m *MarvelClient) weekRange(t time.Time) (time.Time, time.Time) {
	for t.Weekday() != time.Sunday {
		t = t.AddDate(0, 0, -1)
	}

	t = t.AddDate(0, 0, -7)

	return t, t.AddDate(0, 0, 6)
}
