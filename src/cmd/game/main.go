package main

import (
	"flag"
	"os"
	"strings"
	"time"

	"github.com/quasilyte/dicewind/assets"
	"github.com/quasilyte/dicewind/src/controls"
	"github.com/quasilyte/dicewind/src/scenes/encounter"
	"github.com/quasilyte/dicewind/src/session"
	"github.com/quasilyte/ge"

	_ "image/png"
)

func main() {
	state := &session.State{}
	flag.StringVar(&state.AddonDir, "addon-dir", "$DICEWIND_DALE/_addon",
		"extra game data directory")
	flag.Parse()

	if strings.Contains(state.AddonDir, "$DICEWIND_DALE") {
		state.AddonDir = strings.ReplaceAll(state.AddonDir, "$DICEWIND_DALE", os.Getenv("DICEWIND_DALE"))
	}

	ctx := ge.NewContext()
	ctx.Rand.SetSeed(time.Now().Unix())
	ctx.GameName = "dicewind_dale"
	ctx.WindowTitle = "Dicewind Dale"
	ctx.WindowWidth = 1920
	ctx.WindowHeight = 1080
	ctx.FullScreen = true

	assets.Register(ctx)
	controls.BindKeymap(ctx, state)

	// m, err := dataldr.LoadModule(filepath.Join(state.AddonDir, "crypt"))
	// if err != nil {
	// 	panic(err)
	// }
	// level := dungeon.GenerateLevel(m)
	// if err := ge.RunGame(ctx, dungeon.NewController(state, level)); err != nil {
	// 	panic(err)
	// }

	if err := ge.RunGame(ctx, encounter.NewController(state)); err != nil {
		panic(err)
	}
}
