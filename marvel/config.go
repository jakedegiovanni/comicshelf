package marvel

import (
	"github.com/jakedegiovanni/comicshelf/internal/comicclient"
)

type Config struct {
	Client        comicclient.Config `mapstructure:"client"`
	DateLayout    string             `mapstructure:"date_layout"`
	ReleaseOffset int                `mapstructure:"release_offset"`
}
