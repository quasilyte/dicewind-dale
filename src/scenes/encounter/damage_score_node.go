package encounter

import (
	"image/color"
	"strconv"

	"github.com/quasilyte/dicewind/assets"
	"github.com/quasilyte/ge"
)

type damageScoreNode struct {
	damage    int
	pos       ge.Pos
	bg        *ge.Rect
	label     *ge.Label
	textColor color.RGBA
}

var (
	normalDamageColor = ge.RGB(0xe11d51)
	poisonDamageColor = ge.RGB(0x63ef00)
)

func newDamageScoreNode(damage int, pos ge.Pos, c color.RGBA) *damageScoreNode {
	return &damageScoreNode{damage: damage, pos: pos, textColor: c}
}

func (s *damageScoreNode) Init(scene *ge.Scene) {
	s.bg = ge.NewRect(scene.Context(), 64, 32)
	s.bg.Pos = s.pos
	s.bg.FillColorScale.SetRGBA(0x10, 0x10, 0x10, 0xaa)
	scene.AddGraphics(s.bg)

	s.label = scene.NewLabel(assets.FontSmall)
	s.label.Text = strconv.Itoa(s.damage)
	s.label.AlignHorizontal = ge.AlignHorizontalCenter
	s.label.AlignVertical = ge.AlignVerticalCenter
	s.label.Width = s.bg.Width
	s.label.Height = s.bg.Height
	s.label.Pos = s.bg.AnchorPos().WithOffset(0, -4)
	if s.damage != 0 {
		s.label.ColorScale.SetColor(s.textColor)
	}
	scene.AddGraphics(s.label)
}

func (s *damageScoreNode) IsDisposed() bool {
	return s.label.IsDisposed()
}

func (s *damageScoreNode) Update(delta float64) {
	s.label.Pos.Offset.Y -= 32 * delta
	s.bg.Pos.Offset.Y -= 32 * delta

	s.label.ColorScale.A -= float32(delta * 0.4)
	s.bg.FillColorScale.A -= float32(delta * 0.4)
	if s.label.ColorScale.A < 0.1 {
		s.bg.Dispose()
		s.label.Dispose()
	}
}
