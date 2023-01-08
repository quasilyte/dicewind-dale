package encounter

import (
	"github.com/quasilyte/dicewind/assets"
	"github.com/quasilyte/dicewind/src/battle"
	"github.com/quasilyte/dicewind/src/gameui"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/gesignal"
	"github.com/quasilyte/ge/resource"
	"github.com/quasilyte/gmath"
)

type moveRunner struct {
	event *battle.UnitMoveEvent

	sourceLabel *gameui.Label

	image resource.ImageID

	fromPos gmath.Vec
	toPos   gmath.Vec

	progress float64
	pos      gmath.Vec
	sprite   *ge.Sprite

	EventCompleted gesignal.Event[gesignal.Void]
}

func newMoveRunner(e *battle.UnitMoveEvent, fromPos, toPos gmath.Vec, image resource.ImageID) *moveRunner {
	return &moveRunner{
		event: e,
		image: image,

		fromPos: fromPos,
		toPos:   toPos,
		pos:     fromPos,
	}
}

func (r *moveRunner) Init(scene *ge.Scene) {
	r.sourceLabel = gameui.NewLabel(gameui.LabelConfig{
		Width:  196,
		Height: 64,
		Pos:    ge.Pos{Base: &r.event.Unit.Pos},
		Font:   assets.FontTiny,
		Text:   "Move",
	})
	scene.AddObject(r.sourceLabel)

	r.sprite = scene.NewSprite(r.image)
	r.sprite.Pos.Base = &r.pos
	scene.AddGraphics(r.sprite)
}

func (r *moveRunner) IsDisposed() bool { return r.sourceLabel.IsDisposed() }

func (r *moveRunner) Dispose() {
	r.sourceLabel.Dispose()
	r.sprite.Dispose()
}

func (r *moveRunner) Update(delta float64) {
	r.progress += delta * 3
	if r.progress >= 1 {
		r.EventCompleted.Emit(gesignal.Void{})
		r.Dispose()
	} else {
		gravity := gmath.Vec{Y: 96 * 6}
		r.pos = r.fromPos.CubicInterpolate(r.fromPos.Add(gravity), r.toPos, r.toPos.Add(gravity), r.progress)
	}
}
