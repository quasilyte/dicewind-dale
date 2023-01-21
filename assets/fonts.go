package assets

import (
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/resource"
)

func registerFonts(ctx *ge.Context) {
	fontResources := map[resource.FontID]resource.Font{
		// FontVeryTiny: {Path: "font/gondola.ttf", Size: 16},
		FontTiny:  {Path: "font/gondola.ttf", Size: 16},
		FontSmall: {Path: "font/gondola.ttf", Size: 22},
	}
	for id, res := range fontResources {
		ctx.Loader.FontRegistry.Set(id, res)
		ctx.Loader.PreloadFont(id)
	}
}

const (
	FontNone resource.FontID = iota
	FontVeryTiny
	FontTiny
	FontSmall
)
