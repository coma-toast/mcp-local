package webui

import (
	"embed"
	"io/fs"
)

//go:embed all:dist
var embeddedDist embed.FS

func defaultAssets() fs.FS {
	sub, err := fs.Sub(embeddedDist, "dist")
	if err != nil {
		return embeddedDist
	}
	return sub
}
