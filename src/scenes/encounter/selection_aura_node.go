package encounter

import (
	assets "github.com/quasilyte/dicewind/assets"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmath"
)

type selectionAuraNode struct {
	rotation gmath.Rad
	sprite   *ge.Sprite

	Pos gmath.Vec
}

func newSelectionAuraNode() *selectionAuraNode {
	return &selectionAuraNode{}
}

func (a *selectionAuraNode) Init(scene *ge.Scene) {
	a.sprite = scene.NewSprite(assets.ImageSelectionAura)
	a.sprite.Rotation = &a.rotation
	a.sprite.Pos.Base = &a.Pos
	a.sprite.Visible = false
	scene.AddGraphics(a.sprite)
}

func (a *selectionAuraNode) SetVisibility(v bool) {
	a.sprite.Visible = v
}

func (a *selectionAuraNode) IsDisposed() bool { return false }

func (a *selectionAuraNode) Update(delta float64) {
	a.rotation += gmath.Rad(delta)
}
