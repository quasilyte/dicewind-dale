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

func (p *botPlayer) attackFront(u *battle.Unit, reachBack bool) (ruleset.TilePos, bool) {
	attackPos := u.TilePos.OtherAlliance().FrontRow()
	if p.board.Tiles[attackPos.GlobalIndex()].Unit != nil {
		return attackPos, true
	}
	// Can we attack the back row?
	if (reachBack || !p.calc.HasMeleeUnits(u.EnemyAlliance())) && p.board.Tiles[attackPos.BackRow().GlobalIndex()].Unit != nil {
		return attackPos.BackRow(), true
	}
	return ruleset.TilePos{}, false
}

func (p *botPlayer) adjustAction(u *battle.Unit, a ruleset.Action) ruleset.Action {
	switch a.Kind {
	case ruleset.ActionGuard:
		// Nothing to do.

	case ruleset.ActionSkill:
		skill := u.Skill(a.SubKind)
		switch skill.TargetKind {
		case ruleset.TargetEnemyAny:
			attackPos, ok := p.findAttackTarget(u, ruleset.ReachRanged)
			if ok {
				a.Pos = attackPos
				break
			}
			return p.adjustAction(u, ruleset.Action{Kind: ruleset.ActionAttack})

		case ruleset.TargetEnemySpear:
			attackPos, ok := p.findAttackTarget(u, ruleset.ReachRangedFront)
			if ok {
				a.Pos = attackPos
				break
			}
			return p.adjustAction(u, ruleset.Action{Kind: ruleset.ActionAttack})

		default:
			panic("unimplemented")
		}

	case ruleset.ActionMove:
		// 1. Step to the back row.
		if u.HP < u.Monster.HP && !u.TilePos.IsBackRow() {
			tile := p.board.Tiles[u.TilePos.BackRow().GlobalIndex()]
			if tile.Unit == nil {
				a.Pos = u.TilePos.BackRow()
				break
			}
		}
		// 2. Move to the right.
		if p.board.Tiles[u.TilePos.RightPos().GlobalIndex()].Unit == nil {
			a.Pos = u.TilePos.RightPos()
			break
		}
		// 3. Move to the left.
		if p.board.Tiles[u.TilePos.LeftPos().GlobalIndex()].Unit == nil {
			a.Pos = u.TilePos.LeftPos()
			break
		}
		// Otherwise fallback to attack.
		return p.adjustAction(u, ruleset.Action{Kind: ruleset.ActionAttack})

	case ruleset.ActionAttack:
		switch u.Monster.Reach {
		case ruleset.ReachRanged:
			attackPos, ok := p.findAttackTarget(u, u.Monster.Reach)
			if ok {
				a.Pos = attackPos
				break
			}
			return p.adjustAction(u, ruleset.Action{Kind: ruleset.ActionGuard})

		case ruleset.ReachMeleeFront:
			attackPos, ok := p.findAttackTarget(u, u.Monster.Reach)
			if ok {
				a.Pos = attackPos
				break
			}
			// There are no valid targets for the current pos,
			// see if we can move somewhere to reach any other target.
			for pos := (ruleset.TilePos{Alliance: u.Alliance}); pos.Index < 6; pos.Index++ {
				if pos.Index == u.TilePos.Index {
					continue
				}
				if p.board.Tiles[pos.GlobalIndex()].Unit != nil {
					continue // Pos already occupied
				}
				opposingPos := pos.OtherAlliance()
				ok := (p.board.Tiles[opposingPos.GlobalIndex()].Unit != nil && p.calc.CanAttackPos(u, pos, opposingPos)) ||
					(p.board.Tiles[opposingPos.OtherRow().GlobalIndex()].Unit != nil && p.calc.CanAttackPos(u, pos, opposingPos.OtherRow()))
				if ok {
					return ruleset.Action{Kind: ruleset.ActionMove, Pos: pos}
				}
			}
			return p.adjustAction(u, ruleset.Action{Kind: ruleset.ActionGuard})

		case ruleset.ReachMelee:
			attackPos, ok := p.findAttackTarget(u, u.Monster.Reach)
			if ok {
				a.Pos = attackPos
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

func (p *botPlayer) findAttackTarget(u *battle.Unit, reach ruleset.AttackReach) (ruleset.TilePos, bool) {
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
		attackPos.Index = uint8(p.dice.Roll1d6("bot", u.Name(), "attack target"))
		if p.board.Tiles[attackPos.GlobalIndex()].Unit != nil {
			return attackPos, true
		}
		// Otherwise pick a target using enumeration.
		for index := uint8(0); index < 6; index++ {
			attackPos.Index = index
			if p.board.Tiles[attackPos.GlobalIndex()].Unit != nil {
				return attackPos, true
			}
		}
		return ruleset.TilePos{}, false

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
		attackPos.Index = uint8(p.dice.Roll1d6("bot", u.Name(), "attack target"))
		if p.calc.CanAttackPos(u, u.TilePos, attackPos) && p.board.Tiles[attackPos.GlobalIndex()].Unit != nil {
			return attackPos, true
		}
		// Otherwise pick a target using enumeration.
		for index := uint8(0); index < 6; index++ {
			attackPos.Index = index
			if p.board.Tiles[attackPos.GlobalIndex()].Unit != nil {
				return attackPos, true
			}
		}
		return ruleset.TilePos{}, false

	default:
		panic("unimplemented")
	}
}
