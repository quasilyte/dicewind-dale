package assets

import (
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/resource"
)

func registerShaders(ctx *ge.Context) {
	// Associate shader resources.
	shaderResources := map[resource.ShaderID]resource.ShaderInfo{
		ShaderDissolve: {Path: "shader/dissolve.go"},
	}
	for id, res := range shaderResources {
		ctx.Loader.ShaderRegistry.Set(id, res)
		ctx.Loader.PreloadShader(id)
	}
}

const (
	ShaderDissolve resource.ShaderID = iota
)
