package encounter

import (
	"github.com/quasilyte/dicewind/assets"
	"github.com/quasilyte/dicewind/src/gameui"
	"github.com/quasilyte/dicewind/src/ruleset"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmath"
)

type skillSlotNode struct {
	scene *ge.Scene

	pos gmath.Vec

	key string

	border *ge.Sprite
	sprite *ge.Sprite

	label *gameui.LabelWithShadow

	healthCost []*ge.Sprite
	energyCost []*ge.Sprite

	Disabled bool
	Skill    *ruleset.Skill
	Rect     gmath.Rect
}

func newSkillSlotNode(key string, pos gmath.Vec) *skillSlotNode {
	return &skillSlotNode{
		pos: pos,
		key: key,
	}
}

func (s *skillSlotNode) Init(scene *ge.Scene) {
	s.scene = scene

	s.sprite = ge.NewSprite(scene.Context())
	s.sprite.Pos.Base = &s.pos
	s.sprite.Visible = false
	scene.AddGraphics(s.sprite)

	s.label = gameui.NewLabelWithShadow(gameui.LabelWithShadowConfig{
		Pos: ge.Pos{
			Base:   &s.pos,
			Offset: gmath.Vec{X: (-160 / 2) + 4, Y: 38},
		},
		Width: 160 - 28,
		Font:  assets.FontSmall,
	})
	s.label.Text = s.key
	s.label.Visible = false
	scene.AddObject(s.label)

	s.border = scene.NewSprite(assets.ImageSkillBorder)
	s.border.Visible = false
	s.border.Pos.Base = &s.pos
	scene.AddGraphics(s.border)

	s.energyCost = make([]*ge.Sprite, ruleset.MaxEnergyCost)
	energyCostOffset := gmath.Vec{X: 71, Y: 54}
	for i := range s.energyCost {
		sprite := scene.NewSprite(assets.ImageEnergyCost)
		sprite.Visible = false
		sprite.Pos.Base = &s.pos
		sprite.Pos.Offset = energyCostOffset
		scene.AddGraphics(sprite)
		energyCostOffset = energyCostOffset.Sub(gmath.Vec{Y: 16 + 2})
		s.energyCost[i] = sprite
	}

	s.healthCost = make([]*ge.Sprite, ruleset.MaxEnergyCost)
	healthCostOffset := gmath.Vec{X: -71, Y: 54}
	for i := range s.healthCost {
		sprite := scene.NewSprite(assets.ImageHealthCost)
		sprite.Visible = false
		sprite.Pos.Base = &s.pos
		sprite.Pos.Offset = healthCostOffset
		scene.AddGraphics(sprite)
		healthCostOffset = healthCostOffset.Sub(gmath.Vec{Y: 16 + 2})
		s.healthCost[i] = sprite
	}

	s.Rect.Min = s.border.AnchorPos().Resolve()
	s.Rect.Max = s.Rect.Min.Add(gmath.Vec{X: s.border.ImageWidth(), Y: s.border.ImageHeight()})
}

func (s *skillSlotNode) IsDisposed() bool { return false }

func (s *skillSlotNode) Update(delta float64) {}

func (s *skillSlotNode) SetDisabled(disabled bool) {
	if disabled == s.Disabled {
		return
	}
	s.Disabled = disabled
	if disabled {
		s.label.SetColor(ge.RGB(0xe11d51))
	} else {
		s.label.SetColor(ge.RGB(0xffffff))
	}
}

func (s *skillSlotNode) SetSkill(skill *ruleset.Skill) {
	s.Skill = skill

	if skill != nil {
		s.sprite.SetImage(s.scene.LoadImage(skill.Icon))
	}

	visible := skill != nil
	s.border.Visible = visible
	s.sprite.Visible = visible
	s.label.Visible = visible

	for i, sprite := range s.energyCost {
		sprite.Visible = skill != nil && i < skill.EnergyCost
	}
	for i, sprite := range s.healthCost {
		sprite.Visible = skill != nil && i < skill.HealthCost
	}
}
