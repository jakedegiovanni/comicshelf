package static

import "embed"

//go:embed *
var Files embed.FS

type Content struct {
	Date         string
	PageEndpoint string
	Resp         interface{}
}
