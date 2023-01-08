package battle

import (
	"fmt"

	"github.com/quasilyte/dicewind/src/ruleset"
	"github.com/quasilyte/gmath"
)

type Runner struct {
	calc  *Calculator
	dice  *ruleset.Dice
	board *Board

	unit      *Unit
	turnQueue []*Unit

	events     []Event
	postEvents []Event
}

func NewRunner(calc *Calculator, dice *ruleset.Dice, board *Board) *Runner {
	r := &Runner{
		calc:      calc,
		dice:      dice,
		board:     board,
		events:    make([]Event, 0, 8),
		turnQueue: make([]*Unit, 0, 12),
	}
	return r
}

func (r *Runner) StartRound() *Unit {
	r.events = r.events[:0]

	r.buildTurnQueue()
	r.pickNextUnit()
	return r.unit
}

func (r *Runner) ApplyActions(actions []ruleset.Action) (*Unit, []Event) {
	r.events = r.events[:0]
	r.postEvents = r.postEvents[:0]
	r.applyActions(actions)
	r.pickNextUnit()
	r.checkVictory()

	if len(r.postEvents) != 0 {
		r.events = append(r.events, r.postEvents...)
	}

	return r.unit, r.events
}

func (r *Runner) checkVictory() {
	units := [2]uint8{}
	r.board.WalkUnits(func(u *Unit) bool {
		if u.HP <= 0 {
			r.postEvents = append(r.postEvents, &UnitDefeatedEvent{Unit: u})
			r.board.Tiles[u.TilePos.GlobalIndex()].Unit = nil
			return true
		}
		units[u.Alliance]++
		return true
	})
	victoryAlliance := -1
	switch {
	case units[0] == 0:
		victoryAlliance = 1
	case units[1] == 0:
		victoryAlliance = 0
	}
	if victoryAlliance != -1 {
		r.events = append(r.events, &VictoryEvent{Alliance: victoryAlliance})
	}
}

func (r *Runner) applySkillEffects(u *Unit, skill *ruleset.Skill, pos ruleset.TilePos) {
	target := r.board.Tiles[pos.GlobalIndex()].Unit

	rollBonus := 0

	for _, e := range skill.TargetEffects {
		switch e.Kind {
		case ruleset.EffectRollBonus:
			rollBonus += e.Value.(int)
		case ruleset.EffectPoison:
			target.Poison = gmath.ClampMax(target.Poison+e.Value.(int), ruleset.MaxPoison)
		case ruleset.EffectAttack:
			damage := r.calc.AttackDamage(u, rollBonus)
			r.applyDamage(u, ruleset.SourcePhysical, damage, pos)
		case ruleset.EffectDamage:
			damage := r.calc.SkillDamage(u, skill, e)
			r.applyDamage(u, e.Source, damage, target.TilePos)
		case ruleset.EffectSummonSkeleton:
			r.board.AddUnit(NewMonsterUnit(u.Alliance, ruleset.MonsterByName("Skeleton")), pos)
		case ruleset.EffectPoisonToHealth:
			cured := gmath.ClampMax(e.Value.(int), target.Poison)
			target.HP = gmath.ClampMax(target.HP+cured, target.MaxHP())
			target.Poison -= cured
		default:
			panic("unexpected skill effect kind")
		}
	}
}

func (r *Runner) applyActions(actions []ruleset.Action) {
	u := r.unit

	u.Guarding = false

	for _, a := range actions {
		switch a.Kind {
		case ruleset.ActionSkill:
			skill := u.Skill(a.SubKind)
			r.events = append(r.events, &UnitSkillCastEvent{
				Caster: u,
				Target: a.Pos,
				Skill:  skill,
			})
			u.MP -= skill.EnergyCost
			u.HP -= skill.HealthCost
			r.applySkillEffects(u, skill, a.Pos)

		case ruleset.ActionMove:
			r.events = append(r.events, &UnitMoveEvent{
				Unit: u,
				From: u.TilePos,
				To:   a.Pos,
			})
			r.board.Tiles[u.TilePos.GlobalIndex()].Unit = nil
			u.TilePos = a.Pos
			r.board.Tiles[a.Pos.GlobalIndex()].Unit = u

		case ruleset.ActionAttack:
			target := r.board.Tiles[a.Pos.GlobalIndex()].Unit
			r.events = append(r.events, &UnitAttackEvent{
				Attacker: u,
				Defender: target,
			})
			damage := r.calc.AttackDamage(u, 0)
			r.applyDamage(u, ruleset.SourcePhysical, damage, a.Pos)

		case ruleset.ActionGuard:
			u.Guarding = true
			r.events = append(r.events, &UnitGuardEvent{Unit: u})

		default:
			panic(fmt.Sprintf("unhandled %s action", a.Kind))
		}
	}

	if u.Poison > 0 {
		u.Poison--
		u.HP--
		r.postEvents = append(r.postEvents, &UnitDamagedEvent{
			Unit:     u,
			Damage:   1,
			IsPoison: true,
		})
	}
}

func (r *Runner) applyDamage(u *Unit, damageKind ruleset.EffectSource, damage int, pos ruleset.TilePos) {
	target := r.board.Tiles[pos.GlobalIndex()].Unit
	if target.Guarding && damageKind == ruleset.SourcePhysical {
		damage -= 1
	}
	if damage < 0 {
		damage = 0
	}
	target.HP -= damage
	r.events = append(r.events, &UnitDamagedEvent{
		Unit:   target,
		Damage: damage,
	})
}

func (r *Runner) buildTurnQueue() {
	r.turnQueue = r.turnQueue[:0]

	unitIndexMapping := [12][2]uint8{
		0: {0, 0}, 1: {1, 0},
		2: {0, 1}, 3: {1, 1},
		4: {0, 2}, 5: {1, 2},

		6: {0, 3}, 7: {1, 3},
		8: {0, 4}, 9: {1, 4},
		10: {0, 5}, 11: {1, 5},
	}

	for i := len(unitIndexMapping) - 1; i >= 0; i-- {
		indexes := unitIndexMapping[i]
		alliance := indexes[0]
		pos := indexes[1]
		tile := r.board.GetTile(alliance, pos)
		if tile.Unit != nil {
			r.turnQueue = append(r.turnQueue, tile.Unit)
		}
	}
}

func (r *Runner) pickNextUnit() {
	if len(r.turnQueue) == 0 {
		r.unit = nil
		return
	}

	r.unit = r.turnQueue[len(r.turnQueue)-1]
	r.turnQueue = r.turnQueue[:len(r.turnQueue)-1]

	if r.unit.HP <= 0 {
		r.pickNextUnit()
	}
}
