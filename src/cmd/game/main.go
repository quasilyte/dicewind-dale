package main

import (
	"time"

	assets "github.com/quasilyte/dicewind/assets"
	"github.com/quasilyte/dicewind/src/controls"
	"github.com/quasilyte/dicewind/src/scenes/encounter"
	"github.com/quasilyte/dicewind/src/session"
	"github.com/quasilyte/ge"

	_ "image/png"
)

func main() {
	ctx := ge.NewContext()
	ctx.Rand.SetSeed(time.Now().Unix())
	ctx.GameName = "dicewind_dale"
	ctx.WindowTitle = "Dicewind Dale"
	ctx.WindowWidth = 1920
	ctx.WindowHeight = 1080
	ctx.FullScreen = true

	state := &session.State{}

	assets.Register(ctx)
	controls.BindKeymap(ctx, state)

	if err := ge.RunGame(ctx, encounter.NewController(state)); err != nil {
		panic(err)
	}
}
