package marvel

import (
	"github.com/jakedegiovanni/comicshelf/comicclient"
)

type Config struct {
	Client comicclient.Config `mapstructure:"client"`
}
