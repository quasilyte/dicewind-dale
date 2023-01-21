package gameui

import (
	"image/color"

	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/resource"
)

type LabelWithShadow struct {
	Pos     ge.Pos
	Text    string
	Visible bool

	font  resource.FontID
	width float64

	label       *ge.Label
	labelShadow *ge.Label
}

type LabelWithShadowConfig struct {
	Pos ge.Pos

	Width float64

	Font resource.FontID
}

func NewLabelWithShadow(config LabelWithShadowConfig) *LabelWithShadow {
	return &LabelWithShadow{
		Pos:     config.Pos,
		font:    config.Font,
		width:   config.Width,
		Visible: true,
	}
}

func (l *LabelWithShadow) Init(scene *ge.Scene) {
	l.label = scene.NewLabel(l.font)
	l.label.Pos = l.Pos
	l.label.AlignHorizontal = ge.AlignHorizontalRight
	l.label.Width = l.width

	l.labelShadow = scene.NewLabel(l.font)
	l.labelShadow.Pos = l.label.Pos
	l.labelShadow.AlignHorizontal = ge.AlignHorizontalRight
	l.labelShadow.Width = l.width
	l.labelShadow.Pos.Offset.Y += 1
	l.labelShadow.Pos.Offset.X += 1
	l.labelShadow.ColorScale.SetRGBA(0, 0, 0, 0xff)

	scene.AddGraphics(l.labelShadow)
	scene.AddGraphics(l.label)
}

func (l *LabelWithShadow) IsDisposed() bool {
	return l.label.IsDisposed()
}

func (l *LabelWithShadow) Dispose() {
	l.label.Dispose()
	l.labelShadow.Dispose()
}

func (l *LabelWithShadow) SetColor(c color.RGBA) {
	l.label.ColorScale.SetColor(c)
}

func (l *LabelWithShadow) Update(delta float64) {
	l.label.Visible = l.Visible
	l.labelShadow.Visible = l.Visible
	l.label.Text = l.Text
	l.labelShadow.Text = l.Text
}
