package comicshelf

import "context"

type User struct {
	Id        int      `json:"id"`
	Following Set[int] `json:"following"`
}

type UserService interface {
	Following(ctx context.Context, userId, seriesId int) (bool, error)
	Followed(ctx context.Context, userId int) (Set[int], error)
	Follow(ctx context.Context, userId, seriesId int) error
	Unfollow(ctx context.Context, userId, seriesId int) error
}
