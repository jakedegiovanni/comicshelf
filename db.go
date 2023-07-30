package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

type Db struct {
	file     *os.File
	followed map[string]string
	mu       *sync.Mutex
	quit     chan bool
}

func NewDb(filename string) (*Db, error) {
	var f *os.File
	var followed map[string]string

	if _, err := os.Stat(filename); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}

		f, err = os.Create(filename)
		if err != nil {
			return nil, err
		}

		followed = make(map[string]string)
	} else {
		b, err := os.ReadFile(filename)
		if err != nil {
			return nil, err
		}

		if len(b) > 0 {
			err = json.Unmarshal(b, &followed)
			if err != nil {
				if !errors.Is(err, io.EOF) {
					return nil, err
				}
				followed = make(map[string]string)
			}
		} else {
			followed = make(map[string]string)
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
		quit:     make(chan bool),
	}

	db.timedFlush()
	return db, nil
}

func (d *Db) timedFlush() {
	go func() {
		for {
			timer := time.NewTimer(30 * time.Second)
			select {
			case <-d.quit:
				return
			case <-timer.C:
				d.flush()
			}
		}
	}()
}

func (d *Db) flush() {
	_ = d.file.Truncate(0)
	_, _ = d.file.Seek(0, io.SeekStart)

	err := json.NewEncoder(d.file).Encode(d.followed)
	if err != nil {
		log.Println("db save error", err)
		return
	}
	log.Println("db saved")
}

func (d *Db) Shutdown() {
	log.Println("shutting down db")
	defer d.file.Close()
	close(d.quit)
	d.flush()
}

func (d *Db) Follow(series, name string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.followed[series] = name
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
