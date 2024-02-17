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

var errUserNotFound = errors.New("user not found")

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
		followed[0] = comicshelf.User{Id: 0, Following: make(comicshelf.Set[int])} // todo - onboarding process
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
				followed[0] = comicshelf.User{Id: 0, Following: make(comicshelf.Set[int])} // todo - onboarding process
			}
		} else {
			followed = make(map[int]comicshelf.User)
			followed[0] = comicshelf.User{Id: 0, Following: make(comicshelf.Set[int])} // todo - onboarding process
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

func (d *Db) Following(ctx context.Context, userId, seriesId int) (bool, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	user, err := d.getUser(userId)
	if err != nil {
		return false, err
	}

	return user.Following.Has(seriesId), nil
}

func (d *Db) Followed(ctx context.Context, userId int) (comicshelf.Set[int], error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	user, err := d.getUser(userId)
	if err != nil {
		return comicshelf.Set[int]{}, err
	}

	return user.Following, nil
}

func (d *Db) Follow(ctx context.Context, userId, seriesId int) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	user, err := d.getUser(userId)
	if err != nil {
		return err
	}

	user.Following.Put(seriesId)
	d.followed[userId] = user

	slog.Debug(fmt.Sprintf("%+v", d.followed))
	return nil
}

func (d *Db) Unfollow(ctx context.Context, userId, seriesId int) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	user, err := d.getUser(userId)
	if err != nil {
		return err
	}

	user.Following.Delete(seriesId)
	d.followed[userId] = user
	return nil
}

func (d *Db) getUser(userId int) (comicshelf.User, error) {
	user, ok := d.followed[userId]
	if !ok {
		return comicshelf.User{}, fmt.Errorf("no user with id: %d - %w", userId, errUserNotFound)
	}

	return user, nil
}
