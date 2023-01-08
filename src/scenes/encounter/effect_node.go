package encounter

import (
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/gesignal"
	"github.com/quasilyte/ge/resource"
	"github.com/quasilyte/gmath"
)

type effectNode struct {
	pos   gmath.Vec
	image resource.ImageID
	anim  *ge.Animation

	EventCompleted gesignal.Event[gesignal.Void]
}

func newEffectNode(pos gmath.Vec, image resource.ImageID) *effectNode {
	return &effectNode{
		pos:   pos,
		image: image,
	}
}

func (e *effectNode) Init(scene *ge.Scene) {
	s := scene.NewSprite(e.image)
	s.Pos.Base = &e.pos
	scene.AddGraphics(s)

	e.anim = ge.NewAnimation(s, -1)
	e.anim.SetSecondsPerFrame(0.05)
}

func (e *effectNode) IsDisposed() bool {
	return e.anim.IsDisposed()
}

func (e *effectNode) Dispose() {
	e.anim.Sprite().Dispose()
}

func (e *effectNode) Update(delta float64) {
	if e.anim.Tick(delta) {
		e.EventCompleted.Emit(gesignal.Void{})
		e.Dispose()
		return
	}
}
