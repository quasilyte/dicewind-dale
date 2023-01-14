package gameui

import (
	"image/color"

	"github.com/quasilyte/dicewind/assets"
	"github.com/quasilyte/dicewind/src/battle"
	"github.com/quasilyte/dicewind/src/ruleset"
	"github.com/quasilyte/dicewind/src/style"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/gesignal"
	"github.com/quasilyte/ge/physics"
	"github.com/quasilyte/gmath"
)

type UnitTile struct {
	body physics.Body

	scene *ge.Scene

	tilePos ruleset.TilePos

	unit *battle.Unit

	label       *ge.Label
	labelShadow *ge.Label

	bg         *ge.Sprite
	sprite     *ge.Sprite
	auraSprite *ge.Sprite
	border     *ge.Sprite

	hpLevel      *ge.Sprite
	mpLevel      *ge.Sprite
	poisonTokens []*ge.Sprite

	playingDefeat           bool
	defeatShaderProgression float64

	rect gmath.Rect

	action ruleset.ActionKind

	EventAnimationCompleted gesignal.Event[gesignal.Void]
}

func NewUnitTile(pos gmath.Vec, tilePos ruleset.TilePos) *UnitTile {
	tile := &UnitTile{
		tilePos: tilePos,
	}
	tile.body.Pos = pos
	return tile
}

func (tile *UnitTile) Init(scene *ge.Scene) {
	tile.scene = scene

	tile.bg = scene.NewSprite(assets.ImageUnitCardBg)
	tile.bg.Pos.Base = &tile.body.Pos
	scene.AddGraphics(tile.bg)

	tile.sprite = ge.NewSprite(scene.Context())
	tile.sprite.Pos.Base = &tile.body.Pos
	tile.sprite.Visible = false
	tile.sprite.Shader = scene.NewShader(assets.ShaderDissolve)
	tile.sprite.Shader.Texture1 = scene.LoadImage(assets.ImageNoise)
	tile.sprite.Shader.Enabled = false
	scene.AddGraphics(tile.sprite)

	tile.label = scene.NewLabel(assets.FontTiny)
	tile.label.Visible = false
	tile.label.Pos.Base = &tile.body.Pos
	tile.label.Pos.Offset.X = (-256 / 2) + 8
	tile.label.Pos.Offset.Y = 48
	tile.label.AlignHorizontal = ge.AlignHorizontalRight
	tile.label.Width = 256 - 32

	tile.labelShadow = scene.NewLabel(assets.FontTiny)
	tile.labelShadow.Visible = false
	tile.labelShadow.Pos = tile.label.Pos
	tile.labelShadow.AlignHorizontal = ge.AlignHorizontalRight
	tile.labelShadow.Width = 256 - 32
	tile.labelShadow.Pos.Offset.Y += 1
	tile.labelShadow.Pos.Offset.X += 1
	tile.labelShadow.ColorScale.SetRGBA(0, 0, 0, 0xff)

	scene.AddGraphics(tile.labelShadow)
	scene.AddGraphics(tile.label)

	tile.auraSprite = scene.NewSprite(assets.ImageTileSelectionAura)
	tile.auraSprite.Pos.Base = &tile.body.Pos
	tile.auraSprite.Visible = false
	scene.AddGraphics(tile.auraSprite)

	tile.border = scene.NewSprite(assets.ImageUnitBorder)
	tile.border.Pos.Base = &tile.body.Pos
	scene.AddGraphics(tile.border)

	tile.hpLevel = scene.NewSprite(assets.ImageHealthLevel)
	tile.hpLevel.Pos.Base = &tile.body.Pos
	tile.hpLevel.Pos.Offset = gmath.Vec{X: -117}
	tile.hpLevel.Visible = false
	scene.AddGraphics(tile.hpLevel)

	tile.mpLevel = scene.NewSprite(assets.ImageEnergyLevel)
	tile.mpLevel.Pos.Base = &tile.body.Pos
	tile.mpLevel.Pos.Offset = gmath.Vec{X: 116}
	tile.mpLevel.Visible = false
	scene.AddGraphics(tile.mpLevel)

	tile.poisonTokens = make([]*ge.Sprite, ruleset.MaxPoison)
	for i := range tile.poisonTokens {
		token := scene.NewSprite(assets.ImagePoisonToken)
		token.Pos.Base = &tile.body.Pos
		token.Pos.Offset = gmath.Vec{X: 130 + 10*2, Y: (-16 + 48*2) - float64(i*4)}
		token.Visible = false
		scene.AddGraphics(token)
		tile.poisonTokens[i] = token
	}

	tile.rect.Min = tile.bg.AnchorPos().Resolve()
	tile.rect.Max = tile.rect.Min.Add(gmath.Vec{X: tile.bg.ImageWidth(), Y: tile.bg.ImageHeight()})
}

func (tile *UnitTile) IsDisposed() bool {
	return tile.border.IsDisposed()
}

func (tile *UnitTile) Update(delta float64) {
	if tile.playingDefeat {
		tile.defeatShaderProgression -= delta * 2
		tile.sprite.Shader.SetFloatValue("Time", tile.defeatShaderProgression)
		if tile.defeatShaderProgression <= 0 {
			tile.playingDefeat = false
			tile.sprite.Shader.Enabled = false
			tile.EventAnimationCompleted.Emit(gesignal.Void{})
			tile.SetUnit(nil)
		}
		return
	}
}

func (tile *UnitTile) ContainsPos(pos gmath.Vec) bool {
	return tile.rect.Contains(pos)
}

func (tile *UnitTile) GetPos() gmath.Vec {
	return tile.body.Pos
}

func (tile *UnitTile) GetTilePos() ruleset.TilePos {
	return tile.tilePos
}

func (tile *UnitTile) updateUnit() {
	tile.hpLevel.Pos.Offset.Y = float64(85 - (tile.unit.HP * 10))
	tile.hpLevel.Visible = tile.unit.HP > 0
	tile.mpLevel.Pos.Offset.Y = float64(85 - (tile.unit.MP * 10))
	tile.mpLevel.Visible = tile.unit.MP > 0
	for i, token := range tile.poisonTokens {
		token.Visible = i < tile.unit.Poison
	}
}

func (tile *UnitTile) PlayUnitDefeat() {
	tile.playingDefeat = true
	tile.defeatShaderProgression = 1
	tile.sprite.Shader.Enabled = true

	tile.setUnitVisibility(false)
	tile.sprite.Visible = true
}

func (tile *UnitTile) setUnitVisibility(v bool) {
	tile.sprite.Visible = v
	tile.label.Visible = v
	tile.labelShadow.Visible = v
	tile.hpLevel.Visible = v
	tile.mpLevel.Visible = v
	for i, token := range tile.poisonTokens {
		if tile.unit == nil {
			token.Visible = v
		} else {
			token.Visible = v && i < tile.unit.Poison
		}
	}
}

func (tile *UnitTile) GetAction() ruleset.ActionKind {
	return tile.action
}

func (tile *UnitTile) SetAction(kind ruleset.ActionKind) {
	tile.action = kind
	if kind == ruleset.ActionNone {
		tile.auraSprite.Visible = false
		return
	}
	tile.auraSprite.Visible = true

	tile.auraSprite.FlipHorizontal = tile.scene.Rand().Bool()
	tile.auraSprite.FlipVertical = tile.scene.Rand().Bool()

	var c color.RGBA
	switch kind {
	case ruleset.ActionGuard:
		c = style.TileAuraGuardColor
	case ruleset.ActionAttack:
		c = style.TileAuraAttackColor
	case ruleset.ActionMove:
		c = style.TileAuraMoveColor
	case ruleset.ActionSkill:
		c = style.TileAuraCastColor
	}
	var cs ge.ColorScale
	cs.SetColor(c)
	tile.auraSprite.SetColorScale(cs)
}

func (tile *UnitTile) SetUnit(u *battle.Unit) {
	if tile.unit == u && u != nil {
		tile.updateUnit()
		return
	}

	tile.unit = u

	if u != nil {
		tile.sprite.SetImage(tile.scene.LoadImage(u.CardImage()))
		tile.label.Text = u.Name()
		tile.labelShadow.Text = u.Name()
		tile.setUnitVisibility(true)
		tile.updateUnit()
	} else {
		tile.setUnitVisibility(false)
	}
}
