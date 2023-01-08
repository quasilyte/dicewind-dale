package encounter

import (
	"os"

	assets "github.com/quasilyte/dicewind/assets"
	"github.com/quasilyte/dicewind/src/battle"
	"github.com/quasilyte/dicewind/src/controls"
	"github.com/quasilyte/dicewind/src/ruleset"
	"github.com/quasilyte/dicewind/src/session"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/gesignal"
	"github.com/quasilyte/ge/tuple"
	"github.com/quasilyte/gmath"
)

type Controller struct {
	scene        *ge.Scene
	state        *session.State
	dice         *ruleset.Dice
	calc         *battle.Calculator
	runner       *battle.Runner
	eventsRunner *eventsRunner
	board        *battle.Board

	selection *selectionAuraNode

	human    *humanPlayer
	bot      *botPlayer
	nextUnit *battle.Unit

	running bool

	nodes boardNodes
}

type boardNodes struct {
	tiles [2 * 6]*unitTileNode
}

func NewController(state *session.State) *Controller {
	return &Controller{state: state}
}

func (c *Controller) Init(scene *ge.Scene) {
	c.scene = scene

	warriorClass := ruleset.HeroClassByName("Warrior")
	sorcClass := ruleset.HeroClassByName("Sorcerer")
	warriorHero := &ruleset.Hero{
		Name:      "Alpha",
		CurrentHP: warriorClass.HP,
		CurrentMP: warriorClass.MP,
		Class:     warriorClass,
		Weapon: &ruleset.HeroWeapon{
			Class: ruleset.WeaponByName("Sword"),
		},
		CurrentSkills: []*ruleset.Skill{
			ruleset.SkillByName("True Strike"),
			ruleset.SkillByName("Consume Poison"),
		},
	}
	sorcHero := &ruleset.Hero{
		Name:      "Beta",
		CurrentHP: sorcClass.HP,
		CurrentMP: sorcClass.MP,
		Class:     sorcClass,
		Weapon: &ruleset.HeroWeapon{
			Class: ruleset.WeaponByName("Staff"),
		},
		CurrentSkills: []*ruleset.Skill{
			ruleset.SkillByName("Flame Strike"),
			ruleset.SkillByName("Fireball"),
			ruleset.SkillByName("Hellfire"),
			ruleset.SkillByName("Summon Skeleton"),
		},
	}

	c.board = battle.NewBoard()

	c.board.AddUnit(battle.NewHeroUnit(0, sorcHero), ruleset.TilePos{Alliance: 0, Index: 4})
	c.board.AddUnit(battle.NewHeroUnit(0, warriorHero), ruleset.TilePos{Alliance: 0, Index: 0})

	c.board.AddUnit(battle.NewMonsterUnit(1, ruleset.MonsterByName("Grey Minion Archer")), ruleset.TilePos{Alliance: 1, Index: 3})
	c.board.AddUnit(battle.NewMonsterUnit(1, ruleset.MonsterByName("Darkspawn")), ruleset.TilePos{Alliance: 1, Index: 1})
	c.board.AddUnit(battle.NewMonsterUnit(1, ruleset.MonsterByName("Grey Minion")), ruleset.TilePos{Alliance: 1, Index: 2})
	c.board.AddUnit(battle.NewMonsterUnit(1, ruleset.MonsterByName("Lurking Terror")), ruleset.TilePos{Alliance: 1, Index: 0})

	c.dice = ruleset.NewDice(scene.Rand(), os.Stdout)
	c.calc = battle.NewCalculator(c.dice, c.board)
	c.bot = newBotPlayer(c.calc, c.dice, c.board)
	c.human = newHumanPlayer(c.state.MainInput, c.calc, c.board)
	c.runner = battle.NewRunner(c.calc, c.dice, c.board)
	c.eventsRunner = newEventsRunner(&c.nodes)

	c.human.EventActionsReady.Connect(nil, c.onActionsReady)
	c.bot.EventActionsReady.Connect(nil, c.onActionsReady)
	c.eventsRunner.EventCompleted.Connect(nil, c.onEventsCompleted)

	c.selection = newSelectionAuraNode()
	scene.AddObjectAbove(c.selection, 1)

	c.initUI()

	scene.AddObject(c.human)
	scene.AddObject(c.bot)
	scene.AddObject(c.eventsRunner)

	// var events []battle.Event
	// r := battle.NewRunner(dice, board)
	// stop := false
	// for {
	// 	if stop {
	// 		break
	// 	}
	// 	u := r.StartRound()
	// 	if u == nil {
	// 		break
	// 	}
	// 	for {
	// 		actions := c.bot.GetActions(u)
	// 		u, events = r.ApplyActions(actions)
	// 		for _, e := range events {
	// 			// fmt.Println(">", e.Name())
	// 			_, ok := e.(*battle.VictoryEvent)
	// 			if ok {
	// 				stop = true
	// 				break
	// 			}
	// 		}
	// 		if u == nil {
	// 			break
	// 		}
	// 	}
	// }

	// os.Exit(0)
}

func (c *Controller) Update(delta float64) {
	c.handleInput()
}

func (c *Controller) handleInput() {
	if !c.running && c.state.MainInput.ActionIsJustPressed(controls.ActionDebug) {
		c.startRound()
		return
	}
}

func (c *Controller) startRound() {
	c.running = true

	u := c.runner.StartRound()
	if u == nil {
		return
	}
	c.startNextUnitTurn(u, 0)
}

func (c *Controller) onActionsReady(data tuple.Value2[*battle.Unit, []ruleset.Action]) {
	_, actions := data.Fields()
	u, events := c.runner.ApplyActions(actions)
	c.nextUnit = u
	c.eventsRunner.RunEvents(events)
}

func (c *Controller) onEventsCompleted(gesignal.Void) {
	c.updateUnitTiles()

	if c.nextUnit == nil {
		c.running = false
		c.selection.SetVisibility(false)
		return
	}

	c.startNextUnitTurn(c.nextUnit, 0.5)
}

func (c *Controller) startUnitTurn(u *battle.Unit) {
	c.selection.SetVisibility(u.HumanControlled)
	if u.HumanControlled {
		c.selection.Pos = u.Pos
		c.human.StartTurn(u)
	} else {
		c.bot.StartTurn(u)
	}
}

func (c *Controller) startNextUnitTurn(u *battle.Unit, delay float64) {
	if delay == 0 {
		c.startUnitTurn(u)
	} else {
		c.scene.DelayedCall(delay, func() {
			c.startUnitTurn(u)
		})
	}
}

func (c *Controller) updateUnitTiles() {
	c.board.WalkTiles(func(t *battle.Tile) bool {
		c.nodes.tiles[t.TilePos.GlobalIndex()].SetUnit(t.Unit)
		if t.Unit != nil {
			t.Unit.Pos = c.nodes.tiles[t.TilePos.GlobalIndex()].body.Pos
		}
		return true
	})
}

func (c *Controller) initUI() {
	bg := c.scene.NewSprite(assets.ImageEncounterBg)
	bg.Centered = false
	c.scene.AddGraphicsBelow(bg, 1)

	c.board.WalkTiles(func(t *battle.Tile) bool {
		offset := c.calcUnitPos(t.TilePos)
		n := newUnitTileNode(offset, t.TilePos)
		c.scene.AddObject(n)
		c.nodes.tiles[t.TilePos.GlobalIndex()] = n
		c.board.Tiles[t.TilePos.GlobalIndex()].Pos = n.body.Pos
		return true
	})

	c.updateUnitTiles()
}

func (c *Controller) calcUnitPos(pos ruleset.TilePos) gmath.Vec {
	col := float64(pos.Index)
	row := 0.0
	if pos.IsBackRow() {
		col -= 3
		if pos.Alliance == 1 {
			row = 1
		}
	} else {
		if pos.Alliance == 0 {
			row = 1
		}
	}
	extraOffset := gmath.Vec{}
	if pos.Alliance == 1 {
		extraOffset.Y = (456 + 16) + (col * 32)
	} else {
		extraOffset.Y = -(col * 32)
	}
	offset := gmath.Vec{
		X: 208 + (col * (320 + 32)),
		Y: 190 + (row * (196 + 32)),
	}
	return offset.Add(extraOffset)
}
