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
	Armor  *HeroArmor

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

type HeroArmor struct {
	Class *HeroArmorClass
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

type HeroArmorClass struct {
	Name string

	Level int

	Effects []ItemEffect
}

type ItemEffect struct {
	Kind ItemEffectKind

	Value int
}

type ItemEffectKind int

const (
	ItemEffectNone ItemEffectKind = iota
	// Defensive effects.
	ItemEffectPhysicalDamageReduction
	ItemEffectAttackerRollReduction
)
