package dungeon

import (
	"github.com/quasilyte/dicewind/src/ruleset"
	"github.com/quasilyte/dicewind/src/session"
	"github.com/quasilyte/ge"
)

type Controller struct {
	state *session.State

	level *ruleset.DungeonLevel
}

func NewController(state *session.State, level *ruleset.DungeonLevel) *Controller {
	return &Controller{
		state: state,
		level: level,
	}
}

func (c *Controller) Init(scene *ge.Scene) {}

func (c *Controller) Update(delta float64) {}
