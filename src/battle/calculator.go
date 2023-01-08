package battle

import (
	"github.com/quasilyte/dicewind/src/ruleset"
	"github.com/quasilyte/ge/xslices"
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

func (c *Calculator) AttackDamage(attacker *Unit) int {
	damage := attacker.AttackDamage()[c.dice.Roll1d6("calc", attacker.Name(), "attack damage")]
	switch attacker.WeaponMastery() {
	case ruleset.MasterySword:
		if xslices.Contains(attacker.Masteries(), ruleset.MasterySword) {
			damage++
		}
	}
	return damage
}

func (c *Calculator) CanCastThere(u *Unit, skill *ruleset.Skill, from, pos TilePos) bool {
	if skill.HealthCost >= u.HP || u.MP < skill.EnergyCost {
		return false
	}

	var reach ruleset.AttackReach
	switch skill.TargetKind {
	case ruleset.TargetEnemyAny:
		reach = ruleset.ReachRanged
	case ruleset.TargetEnemyMelee:
		reach = ruleset.ReachMelee
	case ruleset.TargetEnemySpear:
		reach = ruleset.ReachRangedFront
	case ruleset.TargetEmptyAllied:
		return c.board.Tiles[u.Alliance][pos].Unit == nil
	default:
		panic("unimplemented")
	}
	if skill.CanTargetEnemyTile() {
		return c.canAttackPos(reach, u.EnemyAlliance(), from, pos)
	}
	panic("unimplemented")
}

func (c *Calculator) CanAttackPos(u *Unit, from, pos TilePos) bool {
	if !u.IsMonster() {
		return c.canAttackPos(u.Hero.Weapon.Class.Reach, u.EnemyAlliance(), from, pos)
	}
	return c.canAttackPos(u.Monster.Reach, u.EnemyAlliance(), from, pos)
}

func (c *Calculator) HasMeleeUnits(alliance int) bool {
	for i := 0; i < 3; i++ {
		if c.board.Tiles[alliance][i].Unit != nil {
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

func (c *Calculator) pickTrivialTarget(reach ruleset.AttackReach, alliance int) *Unit {
	switch reach {
	case ruleset.ReachMelee, ruleset.ReachRanged:
		numMeleeTargets := 0
		numTargets := 0
		var meleeTarget *Unit
		var target *Unit
		for i := TilePos(0); i < 6; i++ {
			tile := c.board.Tiles[alliance][i]
			if tile.Unit == nil {
				continue
			}
			target = tile.Unit
			numTargets++
			if !i.IsBackRow() {
				meleeTarget = tile.Unit
			}
		}
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

func (c *Calculator) canAttackPos(reach ruleset.AttackReach, alliance int, from, pos TilePos) bool {
	if c.board.Tiles[alliance][pos].Unit == nil {
		return false
	}

	switch reach {
	case ruleset.ReachRanged:
		return true

	case ruleset.ReachMelee:
		if pos.IsBackRow() {
			return !c.HasMeleeUnits(alliance)
		}
		return true

	case ruleset.ReachMeleeFront:
		if pos.IsBackRow() && c.HasMeleeUnits(alliance) {
			return false
		}
		return pos.Col() == from.Col()

	case ruleset.ReachRangedFront:
		return pos.Col() == from.Col()
	}

	panic("unreachable")
}
