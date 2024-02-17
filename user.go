package comicshelf

import "context"

type User struct {
	Id        int            `json:"id"`
	Following map[int]Series `json:"following"`
}

type UserService interface {
	Followed(ctx context.Context, userId int) ([]Series, error)
	Follow(ctx context.Context, userId, seriesId int) error
	Unfollow(ctx context.Context, userId, seriesId int) error
}
