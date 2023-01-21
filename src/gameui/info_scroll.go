package gameui

import (
	"github.com/quasilyte/dicewind/assets"
	"github.com/quasilyte/dicewind/src/style"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmath"
)

type InfoScroll struct {
	pos gmath.Vec

	bg *ge.Sprite

	title   *ge.Label
	rtitle  *ge.Label
	content *ge.Label
}

func NewInfoScroll(pos gmath.Vec) *InfoScroll {
	return &InfoScroll{
		pos: pos,
	}
}

func (s *InfoScroll) Init(scene *ge.Scene) {
	s.bg = scene.NewSprite(assets.ImagePaperBg)
	s.bg.Pos.Base = &s.pos
	s.bg.Centered = false
	scene.AddGraphics(s.bg)

	s.title = scene.NewLabel(assets.FontSmall)
	s.title.Pos.Base = &s.pos
	s.title.Pos.Offset = gmath.Vec{X: 32, Y: 48}
	s.title.ColorScale.SetColor(style.InfoScrollColor)
	scene.AddGraphics(s.title)

	s.rtitle = scene.NewLabel(assets.FontSmall)
	s.rtitle.Pos = s.title.Pos
	s.rtitle.Width = 600
	s.rtitle.AlignHorizontal = ge.AlignHorizontalRight
	s.rtitle.ColorScale.SetColor(style.InfoScrollColor)
	scene.AddGraphics(s.rtitle)

	s.content = scene.NewLabel(assets.FontTiny)
	s.content.Pos.Base = &s.pos
	s.content.Pos.Offset = gmath.Vec{X: 40, Y: 104}
	s.content.ColorScale.SetColor(style.InfoScrollColor)
	scene.AddGraphics(s.content)
}

func (s *InfoScroll) IsDisposed() bool { return s.bg.IsDisposed() }

func (s *InfoScroll) Update(delta float64) {}

func (s *InfoScroll) SetText(title, rtitle, content string) {
	s.title.Text = title
	s.rtitle.Text = rtitle
	s.content.Text = content
}
