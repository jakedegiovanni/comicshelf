package filedb

import (
	"context"
	"encoding/json"
	"errors"
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
	followed map[string]comicshelf.User
	mu       *sync.RWMutex
	quit     chan bool
	logger   *slog.Logger
}

func New(cfg *Config, logger *slog.Logger) (*Db, error) {
	var f *os.File
	var followed map[string]comicshelf.User

	if _, err := os.Stat(cfg.Filename); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}

		f, err = os.Create(cfg.Filename)
		if err != nil {
			return nil, err
		}

		followed = make(map[string]comicshelf.User)
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
				followed = make(map[string]comicshelf.User)
			}
		} else {
			followed = make(map[string]comicshelf.User)
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
		logger:   logger,
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
		d.logger.Error("db save error", slog.String("err", err.Error()))
		return
	}
	d.logger.Debug("db saved")
}

func (d *Db) Shutdown() {
	d.logger.Debug("shutting down db")
	defer d.file.Close()
	close(d.quit)
	d.flush()
}

func (d *Db) Followed(ctx context.Context, id string) ([]comicshelf.Series, error) {
	panic("not implemented") // TODO: Implement
}

func (d *Db) Follow(ctx context.Context, series comicshelf.Series) error {
	panic("not implemented") // TODO: Implement
}

func (d *Db) Unfollow(ctx context.Context, series comicshelf.Series) error {
	panic("not implemented") // TODO: Implement
}
