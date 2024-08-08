package hooks

import (
	"errors"
	"log/slog"
	"reflect"

	"github.com/mitchellh/mapstructure"
)

func SlogLevelHook() mapstructure.DecodeHookFuncType {
	return func(src, target reflect.Type, data interface{}) (interface{}, error) {
		if src.Kind() != reflect.String {
			return data, nil
		}

		if target != reflect.TypeOf(slog.Level(0)) {
			return data, nil
		}

		dat, ok := data.(string)
		if !ok {
			return nil, errors.New("could not cast decode source to string")
		}

		var lvl slog.Level
		err := lvl.UnmarshalText([]byte(dat))
		return lvl, err
	}
}
