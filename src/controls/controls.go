package controls

import (
	"github.com/quasilyte/dicewind/src/session"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/input"
)

const (
	ActionUnknown input.Action = iota

	ActionDebug

	ActionConfirm
)

func BindKeymap(ctx *ge.Context, state *session.State) {
	keymap := input.Keymap{
		ActionDebug: {input.KeySpace},

		ActionConfirm: {input.KeyMouseLeft},
	}

	state.MainInput = ctx.Input.NewHandler(0, keymap)
}
