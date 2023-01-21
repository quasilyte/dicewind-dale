package encounter

import (
	"os"

	"github.com/quasilyte/dicewind/assets"
	"github.com/quasilyte/dicewind/src/battle"
	"github.com/quasilyte/dicewind/src/controls"
	"github.com/quasilyte/dicewind/src/gameui"
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

	infoScroll *gameui.InfoScroll

	human    *humanPlayer
	bot      *botPlayer
	nextUnit *battle.Unit

	running bool

	nodes boardNodes
}

type boardNodes struct {
	tiles [2 * 6]*gameui.UnitTile
}

func NewController(state *session.State) *Controller {
	return &Controller{state: state}
}

func (c *Controller) Init(scene *ge.Scene) {
	c.scene = scene

	c.infoScroll = gameui.NewInfoScroll(gmath.Vec{X: 1104, Y: 28})
	c.scene.AddObject(c.infoScroll)

	warriorHero := &ruleset.Hero{
		Name:      "Alpha",
		CardImage: assets.ImageHeroWarriorCard,
		Traits: []ruleset.HeroTrait{
			ruleset.TraitStartingHealthBonus,
		},
		Weapon: &ruleset.HeroWeapon{
			Class: ruleset.WeaponByName("Sword"),
		},
		Armor: &ruleset.HeroArmor{
			Class: ruleset.ArmorByName("Mercenary Armor"),
		},
		CurrentSkills: []*ruleset.Skill{
			ruleset.SkillByName("True Strike"),
			ruleset.SkillByName("Consume Poison"),
		},
	}
	warriorHero.CurrentHP = warriorHero.MaxHP()
	warriorHero.CurrentMP = warriorHero.MaxMP()
	sorcHero := &ruleset.Hero{
		Name:      "Beta",
		CardImage: assets.ImageHeroSorcererCard,
		Traits: []ruleset.HeroTrait{
			ruleset.TraitStratingEnergyBonus,
		},
		Weapon: &ruleset.HeroWeapon{
			Class: ruleset.WeaponByName("Staff"),
		},
		CurrentSkills: []*ruleset.Skill{
			ruleset.SkillByName("Flame Strike"),
			ruleset.SkillByName("Fireball"),
			ruleset.SkillByName("Firestorm"),
			ruleset.SkillByName("Hellfire"),
		},
	}
	sorcHero.CurrentHP = sorcHero.MaxHP()
	sorcHero.CurrentMP = sorcHero.MaxMP()

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
	c.human = newHumanPlayer(c.state.MainInput, c.calc, c.board, &c.nodes)
	c.runner = battle.NewRunner(c.calc, c.dice, c.board)
	c.eventsRunner = newEventsRunner(c.board, &c.nodes)

	c.human.EventActionsReady.Connect(nil, c.onActionsReady)
	c.bot.EventActionsReady.Connect(nil, c.onActionsReady)
	c.eventsRunner.EventCompleted.Connect(nil, c.onEventsCompleted)

	c.selection = newSelectionAuraNode()
	scene.AddObjectAbove(c.selection, 1)

	c.initUI()

	scene.AddObject(c.human)
	scene.AddObject(c.bot)
	scene.AddObject(c.eventsRunner)
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
			t.Unit.Pos = c.nodes.tiles[t.TilePos.GlobalIndex()].GetPos()
		}
		return true
	})
}

func (c *Controller) initUI() {
	// bg := c.scene.NewSprite(assets.ImageEncounterBg)
	// bg.Centered = false
	// c.scene.AddGraphicsBelow(bg, 1)
	ctx := c.scene.Context()
	bg := c.scene.NewRepeatedSprite(assets.ImageCryphBg, ctx.WindowWidth, ctx.WindowHeight)
	bg.Centered = false
	c.scene.AddGraphicsBelow(bg, 1)

	c.board.WalkTiles(func(t *battle.Tile) bool {
		offset := gameui.CalcUnitTilePos(t.TilePos)
		n := gameui.NewUnitTile(offset, t.TilePos)
		c.scene.AddObject(n)
		c.nodes.tiles[t.TilePos.GlobalIndex()] = n
		c.board.Tiles[t.TilePos.GlobalIndex()].Pos = offset
		return true
	})

	c.updateUnitTiles()
}
