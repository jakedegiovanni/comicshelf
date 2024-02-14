package comicshelf

import "context"

type User struct {
	Id        string   `json:"id"`
	Following []Series `json:"following"`
}

type UserService interface {
	Followed(ctx context.Context, id string) ([]Series, error)
	Follow(ctx context.Context, series Series) error
	Unfollow(ctx context.Context, series Series) error
}
