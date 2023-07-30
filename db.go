package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"sync"
)

type Db struct {
	file     *os.File
	followed map[string]struct{}
	mu       *sync.Mutex
}

func NewDb(filename string) (*Db, error) {
	var f *os.File
	var followed map[string]struct{}

	if _, err := os.Stat(filename); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}

		f, err = os.Create(filename)
		if err != nil {
			return nil, err
		}
	} else {
		b, err := os.ReadFile(filename)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(b, &followed)
		if err != nil {
			if !errors.Is(err, io.EOF) {
				return nil, err
			}
			followed = make(map[string]struct{})
		}

		f, err = os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
		if err != nil {
			return nil, err
		}
	}

	db := &Db{
		file:     f,
		followed: followed,
		mu:       new(sync.Mutex),
	}

	return db, nil
}

func (d *Db) Shutdown() {
	defer d.file.Close()

	err := json.NewEncoder(d.file).Encode(d.followed)
	if err != nil {
		log.Println("db save error", err)
	}
}

func (d *Db) Follow(series string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.followed[series] = struct{}{}
}

func (d *Db) Unfollow(series string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	delete(d.followed, series)
}

func (d *Db) Following(series string) bool {
	_, ok := d.followed[series]
	return ok
}
