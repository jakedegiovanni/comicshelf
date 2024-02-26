package hooks

import (
	"net/url"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	v := viper.New()

	v.Set("url", "http://localhost:8080/hello")

	var c struct {
		Url url.URL `mapstructure:"url"`
	}
	require.Nil(t, v.Unmarshal(&c, viper.DecodeHook(UrlHook())))

	uri, _ := url.Parse("http://localhost:8080/hello")
	assert.Equal(t, *uri, c.Url)
}
