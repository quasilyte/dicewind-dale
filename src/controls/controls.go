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

	ActionSkill1
	ActionSkill2
	ActionSkill3
	ActionSkill4
)

func BindKeymap(ctx *ge.Context, state *session.State) {
	keymap := input.Keymap{
		ActionDebug: {input.KeySpace},

		ActionConfirm: {input.KeyMouseLeft},

		ActionSkill1: {input.KeyQ},
		ActionSkill2: {input.KeyW},
		ActionSkill3: {input.KeyE},
		ActionSkill4: {input.KeyR},
	}

	state.MainInput = ctx.Input.NewHandler(0, keymap)
}
