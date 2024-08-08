package comicclient

import (
	"net/url"
	"time"
)

type Config struct {
	Timeout time.Duration `mapstructure:"timeout"`
	BaseURL *url.URL      `mapstructure:"base_url"`
}
