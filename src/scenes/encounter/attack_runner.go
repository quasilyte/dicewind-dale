package encounter

import (
	"github.com/quasilyte/dicewind/assets"
	"github.com/quasilyte/dicewind/src/battle"
	"github.com/quasilyte/dicewind/src/gameui"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/gesignal"
)

type attackRunner struct {
	scene *ge.Scene

	event *battle.UnitAttackEvent

	sourceLabel *gameui.Label

	delay float64
	fired bool

	EventCompleted gesignal.Event[gesignal.Void]
}

func newAttackRunner(e *battle.UnitAttackEvent) *attackRunner {
	return &attackRunner{event: e}
}

func (r *attackRunner) Init(scene *ge.Scene) {
	r.scene = scene

	r.sourceLabel = gameui.NewLabel(gameui.LabelConfig{
		Width:  196,
		Height: 64,
		Pos:    ge.Pos{Base: &r.event.Attacker.Pos},
		Font:   assets.FontTiny,
		Text:   "Attack",
	})
	scene.AddObject(r.sourceLabel)

	r.delay = 0.5
}

func (r *attackRunner) IsDisposed() bool { return r.sourceLabel.IsDisposed() }

func (r *attackRunner) Dispose() {
	r.sourceLabel.Dispose()
}

func (r *attackRunner) onAnimationCompleted(gesignal.Void) {
	r.EventCompleted.Emit(gesignal.Void{})
	r.Dispose()
}

func (r *attackRunner) Update(delta float64) {
	if r.fired {
		return
	}
	r.delay -= delta
	if r.delay <= 0 {
		r.fired = true
		r.scene.Audio().PlaySound(r.event.Attacker.AttackSound())
		e := newEffectNode(r.event.Defender.Pos, r.event.Attacker.AttackImage())
		e.EventCompleted.Connect(r, r.onAnimationCompleted)
		r.scene.AddObject(e)
	}
}
