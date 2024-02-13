package templates

import (
	"embed"
)

//go:embed *
var Files embed.FS

type View struct {
	Date string
	Resp interface{}
}
