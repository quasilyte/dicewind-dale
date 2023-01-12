package ruleset

import (
	"github.com/quasilyte/ge/resource"
)

type Hero struct {
	Name string

	CurrentHP int
	CurrentMP int

	CurrentSkills []*Skill

	CardImage resource.ImageID

	Traits []HeroTrait

	Weapon *HeroWeapon
	Armor  *HeroArmor
}

func (h *Hero) DamageRange() DamageRange {
	return h.Weapon.Class.Damage
}

func (h *Hero) WeaponKind() WeaponKind {
	return h.Weapon.Class.Kind
}

func (h *Hero) MaxHP() int {
	hp := 3 // Base hero health
	for _, t := range h.Traits {
		switch t {
		case TraitHealthBonus:
			hp++
		case TraitStartingHealthBonus:
			hp += 3
		case TraitStartingHybridBonus:
			hp += 2
		}
	}
	return hp
}

func (h *Hero) MaxMP() int {
	energy := 2 // Base hero energy
	for _, t := range h.Traits {
		switch t {
		case TraitEnergyBonus:
			energy++
		case TraitStratingEnergyBonus:
			energy += 5
		case TraitStartingHybridBonus:
			energy += 3
		}
	}
	return energy
}

type HeroArmor struct {
	Class *HeroArmorClass
}

type HeroWeapon struct {
	Class *HeroWeaponClass
}

type HeroWeaponClass struct {
	Name string

	Level int

	Kind WeaponKind

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
