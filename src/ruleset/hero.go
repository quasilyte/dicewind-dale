package ruleset

import (
	"github.com/quasilyte/ge/resource"
)

type Hero struct {
	Name string

	CurrentHP int
	CurrentMP int

	CurrentSkills []*Skill

	Weapon *HeroWeapon

	Class *HeroClass
}

func (h *Hero) DamageRange() DamageRange {
	return h.Weapon.Class.Damage
}

func (h *Hero) WeaponKind() WeaponKind {
	return h.Weapon.Class.Kind
}

type HeroClass struct {
	Name string

	CardImage resource.ImageID

	HP int
	MP int

	Masteries []MasteryKind
}

type HeroWeapon struct {
	Class *HeroWeaponClass
}

type MasteryKind int

const (
	MasteryNone MasteryKind = iota
	// MasterySword grants +1 damage for basic attacks when using swords.
	MasterySword
	// MasteryStaff grants +2 magic damage for all spells when rolled 3 or more when using staves.
	MasteryStaff
)

type HeroWeaponClass struct {
	Name string

	Level int

	Kind WeaponKind

	Mastery MasteryKind

	Reach AttackReach

	Damage DamageRange
}
