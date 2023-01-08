package gameui

import (
	"image/color"

	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/resource"
)

type Label struct {
	Pos       ge.Pos
	LabelText string

	width  float64
	height float64

	font  resource.FontID
	label *ge.Label
	bg    *ge.Rect
}

type LabelConfig struct {
	Pos    ge.Pos
	Width  float64
	Height float64
	Font   resource.FontID
	Text   string
}

func NewLabel(config LabelConfig) *Label {
	return &Label{
		Pos:       config.Pos,
		LabelText: config.Text,
		width:     config.Width,
		height:    config.Height,
		font:      config.Font,
	}
}

func (l *Label) Init(scene *ge.Scene) {
	l.bg = ge.NewRect(scene.Context(), l.width, l.height)
	l.bg.Pos = l.Pos
	l.bg.FillColorScale.SetRGBA(0x10, 0x10, 0x10, 0xaa)
	scene.AddGraphics(l.bg)

	l.label = scene.NewLabel(l.font)
	l.label.AlignHorizontal = ge.AlignHorizontalCenter
	l.label.AlignVertical = ge.AlignVerticalCenter
	l.label.Width = l.bg.Width
	l.label.Height = l.bg.Height
	l.label.Pos = l.bg.AnchorPos().WithOffset(0, -4)
	l.label.Text = l.LabelText
	scene.AddGraphics(l.label)
}

func (l *Label) Dispose() {
	l.label.Dispose()
	l.bg.Dispose()
}

func (l *Label) IsDisposed() bool {
	return l.label.IsDisposed()
}

func (l *Label) SetTextColor(c color.RGBA) {
	l.label.ColorScale.SetColor(c)
}

func (l *Label) SetBgColor(c color.RGBA) {
	l.bg.FillColorScale.SetColor(c)
}

func (l *Label) Update(delta float64) {
	l.bg.Pos = l.Pos
	l.label.Pos = l.bg.AnchorPos().WithOffset(0, -4)
	l.label.Text = l.LabelText
}
