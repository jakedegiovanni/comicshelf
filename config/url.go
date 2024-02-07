package config

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"

	"github.com/mitchellh/mapstructure"
)

func UrlHook() mapstructure.DecodeHookFuncType {
	return func(src, target reflect.Type, data interface{}) (interface{}, error) {
		if src.Kind() != reflect.String {
			return data, nil
		}

		if target != reflect.TypeOf(url.URL{}) {
			return data, nil
		}

		dat, ok := data.(string)
		if !ok {
			return nil, errors.New("could not cast decode source to string")
		}

		// todo: empty url

		u, err := url.Parse(dat)
		if err != nil {
			return nil, fmt.Errorf("could not parse url for config decoding: %w", err)
		}

		return *u, nil
	}
}
