package templates

import (
	"embed"
)

//go:embed *
var Files embed.FS

type View[T any] struct {
	Date  string
	Title string
	Resp  T
}
