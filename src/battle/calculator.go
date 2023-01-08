package battle

import (
	"github.com/quasilyte/dicewind/src/ruleset"
	"github.com/quasilyte/ge/xslices"
	"github.com/quasilyte/gmath"
)

type Calculator struct {
	dice  *ruleset.Dice
	board *Board
}

func NewCalculator(dice *ruleset.Dice, board *Board) *Calculator {
	return &Calculator{dice: dice, board: board}
}

func (c *Calculator) SkillDamage(caster *Unit, skill *ruleset.Skill, effect ruleset.Effect) int {
	roll := c.dice.Roll1d6("runner", caster.Name(), "skill effect damage")
	damage := effect.Value.(ruleset.DamageRange)[roll]
	switch effect.Source {
	case ruleset.SourceMagical:
		if caster.WeaponMastery() == ruleset.MasteryStaff && xslices.Contains(caster.Masteries(), ruleset.MasteryStaff) {
			if roll >= 2 {
				damage += 2
			}
		}
	}
	return damage
}

func (c *Calculator) AttackDamage(attacker *Unit, rollBonus int) int {
	roll := c.dice.Roll1d6("calc", attacker.Name(), "attack damage")
	roll = gmath.Clamp(roll+rollBonus, 0, 5)
	damage := attacker.AttackDamage()[roll]
	switch attacker.WeaponMastery() {
	case ruleset.MasterySword:
		if xslices.Contains(attacker.Masteries(), ruleset.MasterySword) {
			damage++
		}
	}
	return damage
}

func (c *Calculator) CanCastThere(u *Unit, skill *ruleset.Skill, from, pos ruleset.TilePos) bool {
	if skill.HealthCost >= u.HP || u.MP < skill.EnergyCost {
		return false
	}

	var reach ruleset.AttackReach
	switch skill.TargetKind {
	case ruleset.TargetSelf:
		return from == pos && pos == u.TilePos
	case ruleset.TargetEnemyAny:
		reach = ruleset.ReachRanged
	case ruleset.TargetEnemyMelee:
		reach = ruleset.ReachMelee
	case ruleset.TargetEnemySpear:
		reach = ruleset.ReachRangedFront
	case ruleset.TargetEmptyAllied:
		return pos.Alliance == uint8(u.Alliance) &&
			c.board.Tiles[pos.GlobalIndex()].Unit == nil
	case ruleset.TargetAttackCandidate:
		reach = u.AttackReach()
	default:
		panic("unimplemented")
	}
	if skill.CanTargetEnemyTile() {
		return c.canAttackPos(reach, from, pos)
	}
	panic("unimplemented")
}

func (c *Calculator) CanAttackPos(u *Unit, from, pos ruleset.TilePos) bool {
	return c.canAttackPos(u.AttackReach(), from, pos)
}

func (c *Calculator) HasMeleeUnits(alliance uint8) bool {
	for i := uint8(0); i < 3; i++ {
		if c.board.GetTile(alliance, i).Unit != nil {
			return true
		}
	}
	return false
}

func (c *Calculator) PickTrivialTarget(u *Unit) *Unit {
	if !u.IsMonster() {
		panic("todo")
	}
	return c.pickTrivialTarget(u.Monster.Reach, u.EnemyAlliance())
}

func (c *Calculator) pickTrivialTarget(reach ruleset.AttackReach, alliance uint8) *Unit {
	switch reach {
	case ruleset.ReachMelee, ruleset.ReachRanged:
		numMeleeTargets := 0
		numTargets := 0
		var meleeTarget *Unit
		var target *Unit
		c.board.WalkTeamUnits(alliance, func(u *Unit) bool {
			target = u
			numTargets++
			if !u.TilePos.IsBackRow() {
				meleeTarget = u
			}
			return true
		})
		if numTargets == 1 {
			return target
		}
		if numMeleeTargets == 1 {
			return meleeTarget
		}
		return nil
	}

	panic("unreachable")
}

func (c *Calculator) canAttackPos(reach ruleset.AttackReach, from, pos ruleset.TilePos) bool {
	if c.board.Tiles[pos.GlobalIndex()].Unit == nil {
		return false
	}

	switch reach {
	case ruleset.ReachRanged:
		return true

	case ruleset.ReachMelee:
		if pos.IsBackRow() {
			return !c.HasMeleeUnits(pos.Alliance)
		}
		return true

	case ruleset.ReachMeleeFront:
		if pos.IsBackRow() && c.HasMeleeUnits(pos.Alliance) {
			return false
		}
		return pos.Col() == from.Col()

	case ruleset.ReachRangedFront:
		return pos.Col() == from.Col()
	}

	panic("unreachable")
}
