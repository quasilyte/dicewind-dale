package battle

import "github.com/quasilyte/dicewind/src/ruleset"

type Event interface {
	Name() string
}

type UnitMoveEvent struct {
	Unit *Unit
	From TilePos
	To   TilePos
}

type UnitSkillCastEvent struct {
	Caster *Unit
	Target TilePos
	Skill  *ruleset.Skill
}

type UnitAttackEvent struct {
	Attacker *Unit
	Defender *Unit
}

type UnitGuardEvent struct {
	Unit *Unit
}

type UnitDamagedEvent struct {
	Unit     *Unit
	IsPoison bool
	Damage   int
}

type UnitDefeatedEvent struct {
	Unit *Unit
}

type VictoryEvent struct {
	Alliance int
}

func (e *UnitMoveEvent) Name() string      { return "UnitMove" }
func (e *UnitSkillCastEvent) Name() string { return "UnitSkillCastEvent" }
func (e *UnitAttackEvent) Name() string    { return "UnitAttackEvent" }
func (e *UnitGuardEvent) Name() string     { return "UnitAGuardEvent" }
func (e *UnitDamagedEvent) Name() string   { return "UnitDamaged" }
func (e *UnitDefeatedEvent) Name() string  { return "UnitDefeated" }
func (e *VictoryEvent) Name() string       { return "Victory" }
