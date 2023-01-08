package ruleset

import "github.com/quasilyte/ge/resource"

type DamageRange [6]int

type ActionRange [6]Action

type ActionKind int

//go:generate stringer -type=ActionKind -trimprefix=Action
const (
	ActionNone ActionKind = iota
	ActionAttack
	ActionGuard
	ActionSkill
	ActionMove
)

type Action struct {
	Kind    ActionKind
	SubKind int
	Pos     TilePos
}

type Race int

const (
	RaceUnknown Race = iota
	RaceGreyMinion
)

type AttackReach int

const (
	ReachNone AttackReach = iota
	ReachMeleeFront
	ReachMelee
	ReachRangedFront
	ReachRanged
)

type WeaponKind int

const (
	WeaponNone WeaponKind = iota
	WeaponClaws
	WeaponSword
	WeaponScimitar
	WeaponBlunt
	WeaponBow
)

type Monster struct {
	Name string

	CardImage resource.ImageID
	Weapon    WeaponKind

	Level  int
	Reward int

	HP int
	MP int

	Damage DamageRange
	Reach  AttackReach

	Action ActionRange

	Skills []*Skill

	Race Race
}
