package templates

import (
	"embed"
)

//go:embed *
var Files embed.FS

type View[T any] struct {
	Date string
	Resp T
}
