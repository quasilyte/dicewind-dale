//go:build ignore
// +build ignore

package main

var Time float

func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
	limit := abs(2*fract(Time/3) - 1)
	level := imageSrc1UnsafeAt(texCoord).x

	if limit-0.1 < level && level < limit {
		alpha := imageSrc0UnsafeAt(texCoord).w
		return vec4(0.2, 0, 0, alpha)
	}

	return step(limit, level) * imageSrc0UnsafeAt(texCoord)
}
