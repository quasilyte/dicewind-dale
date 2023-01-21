package assets

import (
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/resource"
)

func registerRaws(ctx *ge.Context) {
	rawResources := map[resource.RawID]resource.Raw{
		RawDictEn: {Path: "langs/en.txt"},
		RawDictRu: {Path: "langs/ru.txt"},
	}

	for id, res := range rawResources {
		ctx.Loader.RawRegistry.Set(id, res)
		ctx.Loader.PreloadRaw(id)
	}
}

const (
	RawNone resource.RawID = iota

	RawDictEn
	RawDictRu
)
