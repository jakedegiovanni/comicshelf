package comicshelf

import "context"

type User struct {
	Id        string   `json:"id"`
	Following []Series `json:"following"`
}

type UserService interface {
	Followed(ctx context.Context) ([]Series, error)
	Follow(ctx context.Context, seriesId string) error
	Unfollow(ctx context.Context, seriesId string) error
}
