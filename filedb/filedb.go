package filedb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/jakedegiovanni/comicshelf"
)

var _ comicshelf.UserService = (*Db)(nil)

type Db struct {
	file     *os.File
	followed map[int]comicshelf.User
	mu       *sync.RWMutex
	quit     chan bool
}

func New(cfg *Config) (*Db, error) {
	var f *os.File
	var followed map[int]comicshelf.User

	if _, err := os.Stat(cfg.Filename); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}

		f, err = os.Create(cfg.Filename)
		if err != nil {
			return nil, err
		}

		followed = make(map[int]comicshelf.User)
	} else {
		b, err := os.ReadFile(cfg.Filename)
		if err != nil {
			return nil, err
		}

		if len(b) > 0 {
			err = json.Unmarshal(b, &followed)
			if err != nil {
				if !errors.Is(err, io.EOF) {
					return nil, err
				}
				followed = make(map[int]comicshelf.User)
			}
		} else {
			followed = make(map[int]comicshelf.User)
		}

		f, err = os.OpenFile(cfg.Filename, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
		if err != nil {
			return nil, err
		}
	}

	db := &Db{
		file:     f,
		followed: followed,
		mu:       new(sync.RWMutex),
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
		slog.Error("db save error", slog.String("err", err.Error()))
		return
	}
	slog.Debug("db saved")
}

func (d *Db) Shutdown() {
	slog.Debug("shutting down db")
	defer d.file.Close()
	close(d.quit)
	d.flush()
}

func (d *Db) Followed(ctx context.Context, userId int) ([]comicshelf.Series, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	user, err := d.getUser(userId)
	if err != nil {
		return []comicshelf.Series{}, err
	}

	following := make([]comicshelf.Series, 0, len(user.Following))
	for _, series := range user.Following {
		following = append(following, series)
	}

	return following, nil
}

func (d *Db) Follow(ctx context.Context, userId, seriesId int) error {
	d.mu.Lock()
	defer d.mu.Lock()

	user, err := d.getUser(userId)
	if err != nil {
		return err
	}

	if series, ok := user.Following[seriesId]; !ok {
		user.Following[seriesId] = comicshelf.Series{Id: seriesId}
		d.followed[userId] = user
	} else {
		series.Id = seriesId
		user.Following[seriesId] = series
		d.followed[userId] = user
	}

	return nil
}

func (d *Db) Unfollow(ctx context.Context, userId, seriesId int) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	user, err := d.getUser(userId)
	if err != nil {
		return err
	}

	delete(user.Following, seriesId)
	d.followed[userId] = user
	return nil
}

func (d *Db) getUser(userId int) (comicshelf.User, error) {
	user, ok := d.followed[userId]
	if !ok {
		return comicshelf.User{}, fmt.Errorf("no user with id: %d", userId)
	}

	return user, nil
}
