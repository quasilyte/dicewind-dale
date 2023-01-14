package dungeon

import (
	"github.com/quasilyte/dicewind/assets"
	"github.com/quasilyte/dicewind/src/ruleset"
	"github.com/quasilyte/dicewind/src/session"
	"github.com/quasilyte/ge"
)

type Controller struct {
	state *session.State

	scene *ge.Scene

	level *ruleset.DungeonLevel
}

func NewController(state *session.State, level *ruleset.DungeonLevel) *Controller {
	return &Controller{
		state: state,
		level: level,
	}
}

func (c *Controller) Init(scene *ge.Scene) {
	c.scene = scene
	c.initUI()
}

func (c *Controller) Update(delta float64) {}

func (c *Controller) initUI() {
	bg := c.scene.NewSprite(assets.ImageEncounterBg)
	bg.Centered = false
	c.scene.AddGraphicsBelow(bg, 1)
}
