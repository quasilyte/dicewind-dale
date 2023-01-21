package assets

import (
	"embed"
	"io"

	"github.com/quasilyte/ge"
)

func Register(ctx *ge.Context) {
	ctx.Loader.OpenAssetFunc = func(path string) io.ReadCloser {
		f, err := gameAssets.Open("_data/" + path)
		if err != nil {
			ctx.OnCriticalError(err)
		}
		return f
	}

	registerImages(ctx)
	registerFonts(ctx)
	registerAudio(ctx)
	registerShaders(ctx)
	registerRaws(ctx)
}

//go:embed all:_data
var gameAssets embed.FS
