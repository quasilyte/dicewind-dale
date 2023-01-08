package encounter

import (
	"fmt"

	"github.com/quasilyte/dicewind/src/battle"
	"github.com/quasilyte/dicewind/src/ruleset"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/gesignal"
	"github.com/quasilyte/ge/tuple"
)

type botPlayer struct {
	calc  *battle.Calculator
	board *battle.Board
	dice  *ruleset.Dice

	unit *battle.Unit

	active bool

	EventActionsReady gesignal.Event[tuple.Value2[*battle.Unit, []ruleset.Action]]
}

func newBotPlayer(calc *battle.Calculator, dice *ruleset.Dice, board *battle.Board) *botPlayer {
	return &botPlayer{
		calc:  calc,
		board: board,
		dice:  dice,
	}
}

func (p *botPlayer) Init(scene *ge.Scene) {}

func (p *botPlayer) IsDisposed() bool { return false }

func (p *botPlayer) Update(delta float64) {
	if !p.active {
		return
	}

	p.active = false
	actions := p.getActions(p.unit)
	p.EventActionsReady.Emit(tuple.New2(p.unit, actions))
}

func (p *botPlayer) StartTurn(u *battle.Unit) {
	p.active = true
	p.unit = u
}

func (p *botPlayer) getActions(u *battle.Unit) []ruleset.Action {
	if u.Hero != nil || u.Monster == nil {
		panic("bot can only handle monsters for now")
	}

	rolledAction := u.Monster.Action[p.dice.Roll1d6("bot", u.Name(), "choose action")]
	a := p.adjustAction(u, rolledAction)

	return []ruleset.Action{a}
}

func (p *botPlayer) attackFront(u *battle.Unit, reachBack bool) (battle.TilePos, bool) {
	enemyTiles := &p.board.Tiles[u.EnemyAlliance()]
	attackPos := u.TilePos
	if u.TilePos.IsBackRow() {
		attackPos -= 3
	}
	if enemyTiles[attackPos].Unit != nil {
		return attackPos, true
	}
	// Can we attack the back row?
	if (reachBack || !p.calc.HasMeleeUnits(u.EnemyAlliance())) && enemyTiles[attackPos+3].Unit != nil {
		return attackPos + 3, true
	}
	return 0, false
}

func (p *botPlayer) adjustAction(u *battle.Unit, a ruleset.Action) ruleset.Action {
	alliedTiles := &p.board.Tiles[u.Alliance]
	enemyTiles := &p.board.Tiles[u.EnemyAlliance()]

	switch a.Kind {
	case ruleset.ActionGuard:
		// Nothing to do.

	case ruleset.ActionSkill:
		skill := u.Skill(a.SubKind)
		switch skill.TargetKind {
		case ruleset.TargetEnemyAny:
			attackPos, ok := p.findAttackTarget(u, ruleset.ReachRanged)
			if ok {
				a.Value = int(attackPos)
				break
			}
			return p.adjustAction(u, ruleset.Action{Kind: ruleset.ActionAttack})

		case ruleset.TargetEnemySpear:
			attackPos, ok := p.findAttackTarget(u, ruleset.ReachRangedFront)
			if ok {
				a.Value = int(attackPos)
				break
			}
			return p.adjustAction(u, ruleset.Action{Kind: ruleset.ActionAttack})

		default:
			panic("unimplemented")
		}

	case ruleset.ActionMove:
		// 1. Step to the back row.
		if u.HP < u.Monster.HP && !u.TilePos.IsBackRow() {
			tile := alliedTiles[u.TilePos+3]
			if tile.Unit == nil {
				a.Value = int(u.TilePos + 3)
				break
			}
		}
		// 2. Move to the right.
		if alliedTiles[u.TilePos.RightPos()].Unit == nil {
			a.Value = int(u.TilePos.RightPos())
			break
		}
		// 3. Move to the left.
		if alliedTiles[u.TilePos.LeftPos()].Unit == nil {
			a.Value = int(u.TilePos.LeftPos())
			break
		}
		// Otherwise fallback to attack.
		return p.adjustAction(u, ruleset.Action{Kind: ruleset.ActionAttack})

	case ruleset.ActionAttack:
		switch u.Monster.Reach {
		case ruleset.ReachRanged:
			attackPos, ok := p.findAttackTarget(u, u.Monster.Reach)
			if ok {
				a.Value = int(attackPos)
				break
			}
			return p.adjustAction(u, ruleset.Action{Kind: ruleset.ActionGuard})

		case ruleset.ReachMeleeFront:
			attackPos, ok := p.findAttackTarget(u, u.Monster.Reach)
			if ok {
				a.Value = int(attackPos)
				break
			}
			// There are no valid targets for the current pos,
			// see if we can move somewhere to reach any other target.
			for pos := battle.TilePos(0); pos < 6; pos++ {
				if pos == u.TilePos {
					continue
				}
				if alliedTiles[pos].Unit != nil {
					continue // Pos already occupied
				}
				if enemyTiles[pos].Unit != nil && p.calc.CanAttackPos(u, pos, pos) {
					return ruleset.Action{Kind: ruleset.ActionMove, Value: int(pos)}
				}
				if enemyTiles[pos.OtherRow()].Unit != nil && p.calc.CanAttackPos(u, pos, pos.OtherRow()) {
					return ruleset.Action{Kind: ruleset.ActionMove, Value: int(pos)}
				}
			}
			return p.adjustAction(u, ruleset.Action{Kind: ruleset.ActionGuard})

		case ruleset.ReachMelee:
			attackPos, ok := p.findAttackTarget(u, u.Monster.Reach)
			if ok {
				a.Value = int(attackPos)
				break
			}
			return p.adjustAction(u, ruleset.Action{Kind: ruleset.ActionGuard})

		default:
			panic("unimplemented")
		}

	default:
		panic(fmt.Sprintf("unexpected %s action", a.Kind))
	}

	return a
}

func (p *botPlayer) findAttackTarget(u *battle.Unit, reach ruleset.AttackReach) (battle.TilePos, bool) {
	enemyTiles := &p.board.Tiles[u.EnemyAlliance()]

	switch reach {
	case ruleset.ReachRanged:
		attackPos, ok := p.attackFront(u, true)
		if ok {
			return attackPos, true
		}
		target := p.calc.PickTrivialTarget(u)
		if target != nil {
			return target.TilePos, true
		}
		// Do a roll to choose a target.
		attackPos = battle.TilePos(p.dice.Roll1d6("bot", u.Name(), "attack target"))
		if enemyTiles[attackPos].Unit != nil {
			return attackPos, true
		}
		// Otherwise pick a target using enumeration.
		for i := battle.TilePos(0); i < 6; i++ {
			if enemyTiles[i].Unit != nil {
				return i, true
			}
		}
		return 0, false

	case ruleset.ReachRangedFront:
		return p.attackFront(u, true)

	case ruleset.ReachMeleeFront:
		return p.attackFront(u, false)

	case ruleset.ReachMelee:
		attackPos, ok := p.attackFront(u, false)
		if ok {
			return attackPos, true
		}
		target := p.calc.PickTrivialTarget(u)
		if target != nil {
			return target.TilePos, true
		}
		// Do a roll to choose a target.
		attackPos = battle.TilePos(p.dice.Roll1d6("bot", u.Name(), "attack target"))
		if p.calc.CanAttackPos(u, u.TilePos, attackPos) && enemyTiles[attackPos].Unit != nil {
			return attackPos, true
		}
		// Otherwise pick a target using enumeration.
		for i := battle.TilePos(0); i < 6; i++ {
			if enemyTiles[i].Unit != nil {
				return i, true
			}
		}
		return 0, false

	default:
		panic("unimplemented")
	}
}
