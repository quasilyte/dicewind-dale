package encounter

import (
	assets "github.com/quasilyte/dicewind/assets"
	"github.com/quasilyte/dicewind/src/battle"
	"github.com/quasilyte/dicewind/src/ruleset"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/gesignal"
	"github.com/quasilyte/ge/physics"
	"github.com/quasilyte/gmath"
)

type unitTileNode struct {
	body physics.Body

	scene *ge.Scene

	alliance int
	tilePos  battle.TilePos

	unit *battle.Unit

	label       *ge.Label
	labelShadow *ge.Label

	bg     *ge.Sprite
	sprite *ge.Sprite
	border *ge.Sprite

	hpLevel      *ge.Sprite
	mpLevel      *ge.Sprite
	poisonTokens []*ge.Sprite

	playingDefeat           bool
	defeatShaderProgression float64

	EventAnimationCompleted gesignal.Event[gesignal.Void]
}

func newUnitTileNode(pos gmath.Vec, alliance int, tilePos battle.TilePos) *unitTileNode {
	n := &unitTileNode{
		alliance: alliance,
		tilePos:  tilePos,
	}
	n.body.Pos = pos
	return n
}

func (n *unitTileNode) Init(scene *ge.Scene) {
	n.scene = scene

	n.bg = scene.NewSprite(assets.ImageUnitCardBg)
	n.bg.Pos.Base = &n.body.Pos
	scene.AddGraphics(n.bg)

	n.sprite = ge.NewSprite(scene.Context())
	n.sprite.Pos.Base = &n.body.Pos
	n.sprite.Visible = false
	n.sprite.Shader = scene.NewShader(assets.ShaderDissolve)
	n.sprite.Shader.Texture1 = scene.LoadImage(assets.ImageNoise)
	n.sprite.Shader.Enabled = false
	scene.AddGraphics(n.sprite)

	n.label = scene.NewLabel(assets.FontTiny)
	n.label.Visible = false
	n.label.Pos.Base = &n.body.Pos
	n.label.Pos.Offset.X = (-256 / 2) + 8
	n.label.Pos.Offset.Y = 48
	n.label.AlignHorizontal = ge.AlignHorizontalRight
	n.label.Width = 256 - 32

	n.labelShadow = scene.NewLabel(assets.FontTiny)
	n.labelShadow.Visible = false
	n.labelShadow.Pos = n.label.Pos
	n.labelShadow.AlignHorizontal = ge.AlignHorizontalRight
	n.labelShadow.Width = 256 - 32
	n.labelShadow.Pos.Offset.Y += 1
	n.labelShadow.Pos.Offset.X += 1
	n.labelShadow.ColorScale.SetRGBA(0, 0, 0, 0xff)

	scene.AddGraphics(n.labelShadow)
	scene.AddGraphics(n.label)

	n.border = scene.NewSprite(assets.ImageUnitBorder)
	n.border.Pos.Base = &n.body.Pos
	scene.AddGraphics(n.border)

	n.hpLevel = scene.NewSprite(assets.ImageHealthLevel)
	n.hpLevel.Pos.Base = &n.body.Pos
	n.hpLevel.Pos.Offset = gmath.Vec{X: -117}
	n.hpLevel.Visible = false
	scene.AddGraphics(n.hpLevel)

	n.mpLevel = scene.NewSprite(assets.ImageEnergyLevel)
	n.mpLevel.Pos.Base = &n.body.Pos
	n.mpLevel.Pos.Offset = gmath.Vec{X: 116}
	n.mpLevel.Visible = false
	scene.AddGraphics(n.mpLevel)

	n.poisonTokens = make([]*ge.Sprite, ruleset.MaxPoison)
	for i := range n.poisonTokens {
		token := scene.NewSprite(assets.ImagePoisonToken)
		token.Pos.Base = &n.body.Pos
		token.Pos.Offset = gmath.Vec{X: 130 + 10*2, Y: (-16 + 48*2) - float64(i*4)}
		token.Visible = false
		scene.AddGraphics(token)
		n.poisonTokens[i] = token
	}
}

func (n *unitTileNode) IsDisposed() bool {
	return n.border.IsDisposed()
}

func (n *unitTileNode) Update(delta float64) {
	if n.playingDefeat {
		n.defeatShaderProgression -= delta * 2
		n.sprite.Shader.SetFloatValue("Time", n.defeatShaderProgression)
		if n.defeatShaderProgression <= 0 {
			n.playingDefeat = false
			n.sprite.Shader.Enabled = false
			n.EventAnimationCompleted.Emit(gesignal.Void{})
			n.SetUnit(nil)
		}
		return
	}
}

func (n *unitTileNode) updateUnit() {
	n.hpLevel.Pos.Offset.Y = float64(85 - (n.unit.HP * 10))
	n.hpLevel.Visible = n.unit.HP > 0
	n.mpLevel.Pos.Offset.Y = float64(85 - (n.unit.MP * 10))
	n.mpLevel.Visible = n.unit.MP > 0
	for i, token := range n.poisonTokens {
		token.Visible = i < n.unit.Poison
	}
}

func (n *unitTileNode) PlayUnitDefeat() {
	n.playingDefeat = true
	n.defeatShaderProgression = 1
	n.sprite.Shader.Enabled = true

	n.setVisibility(false)
	n.sprite.Visible = true
}

func (n *unitTileNode) setVisibility(v bool) {
	n.sprite.Visible = v
	n.label.Visible = v
	n.labelShadow.Visible = v
	n.hpLevel.Visible = v
	n.mpLevel.Visible = v
	for i, token := range n.poisonTokens {
		if n.unit == nil {
			token.Visible = v
		} else {
			token.Visible = v && i < n.unit.Poison
		}
	}
}

func (n *unitTileNode) SetUnit(u *battle.Unit) {
	if n.unit == u && u != nil {
		n.updateUnit()
		return
	}

	n.unit = u

	if u != nil {
		n.sprite.SetImage(n.scene.LoadImage(u.CardImage()))
		n.label.Text = u.Name()
		n.labelShadow.Text = u.Name()
		n.setVisibility(true)
		n.updateUnit()
	} else {
		n.setVisibility(false)
	}
}
