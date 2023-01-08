package assets

import (
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/resource"
)

func registerFonts(ctx *ge.Context) {
	fontResources := map[resource.FontID]resource.Font{
		FontVeryTiny: {Path: "font/Fiolex_Mephisto.otf", Size: 14},
		FontTiny:     {Path: "font/Fiolex_Mephisto.otf", Size: 18},
		FontSmall:    {Path: "font/Fiolex_Mephisto.otf", Size: 24},
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
