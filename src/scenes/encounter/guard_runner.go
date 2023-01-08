package encounter

import (
	"github.com/quasilyte/dicewind/assets"
	"github.com/quasilyte/dicewind/src/battle"
	"github.com/quasilyte/dicewind/src/gameui"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/gesignal"
)

type guardRunner struct {
	event *battle.UnitGuardEvent

	sourceLabel *gameui.Label

	delay float64

	EventCompleted gesignal.Event[gesignal.Void]
}

func newGuardRunner(e *battle.UnitGuardEvent) *guardRunner {
	return &guardRunner{event: e}
}

func (r *guardRunner) Init(scene *ge.Scene) {
	r.sourceLabel = gameui.NewLabel(gameui.LabelConfig{
		Width:  196,
		Height: 64,
		Pos:    ge.Pos{Base: &r.event.Unit.Pos},
		Font:   assets.FontTiny,
		Text:   "Guard",
	})
	scene.AddObject(r.sourceLabel)

	r.delay = 0.6
}

func (r *guardRunner) IsDisposed() bool { return r.sourceLabel.IsDisposed() }

func (r *guardRunner) Dispose() {
	r.sourceLabel.Dispose()
}

func (r *guardRunner) Update(delta float64) {
	r.delay -= delta
	if r.delay <= 0 {
		r.EventCompleted.Emit(gesignal.Void{})
		r.Dispose()
	}
}
