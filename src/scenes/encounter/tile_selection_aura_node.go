package encounter

import (
	"image/color"

	assets "github.com/quasilyte/dicewind/assets"
	"github.com/quasilyte/dicewind/src/ruleset"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmath"
)

var (
	tileColorMove   = color.RGBA{R: 0x0e, G: 0x8d, B: 0xcd, A: 120}
	tileColorGuard  = color.RGBA{R: 0x56, G: 0xb0, B: 0x66, A: 120}
	tileColorAttack = color.RGBA{R: 0xcc, G: 0x31, B: 0x04, A: 120}
	tileColorCast   = color.RGBA{R: 0x99, G: 0x41, B: 0xb5, A: 120}
)

type tileSelectionAuraNode struct {
	TilePos ruleset.TilePos

	pos gmath.Vec

	rand *gmath.Rand

	sprite *ge.Sprite

	Action ruleset.ActionKind
	Rect   gmath.Rect
}

func newTileSelectionAuraNode(tilePos ruleset.TilePos, pos gmath.Vec) *tileSelectionAuraNode {
	return &tileSelectionAuraNode{
		TilePos: tilePos,
		pos:     pos,
	}
}

func (a *tileSelectionAuraNode) Init(scene *ge.Scene) {
	a.rand = scene.Rand()

	a.sprite = scene.NewSprite(assets.ImageTileSelectionAura)
	a.sprite.Pos.Base = &a.pos
	a.sprite.Visible = false
	scene.AddGraphics(a.sprite)

	a.Rect.Min = a.sprite.AnchorPos().Resolve()
	a.Rect.Max = a.Rect.Min.Add(gmath.Vec{X: a.sprite.ImageWidth(), Y: a.sprite.ImageHeight()})
}

func (a *tileSelectionAuraNode) IsDisposed() bool { return false }

func (a *tileSelectionAuraNode) SetVisibility(v bool) {
	a.sprite.Visible = v
}

func (a *tileSelectionAuraNode) SetAction(kind ruleset.ActionKind) {
	a.Action = kind

	a.sprite.FlipHorizontal = a.rand.Bool()
	a.sprite.FlipVertical = a.rand.Bool()

	var c color.RGBA
	switch kind {
	case ruleset.ActionGuard:
		c = tileColorGuard
	case ruleset.ActionAttack:
		c = tileColorAttack
	case ruleset.ActionMove:
		c = tileColorMove
	case ruleset.ActionSkill:
		c = tileColorCast
	}
	var cs ge.ColorScale
	cs.SetColor(c)
	a.sprite.SetColorScale(cs)
}

func (a *tileSelectionAuraNode) Update(delta float64) {}
