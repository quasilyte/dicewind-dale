package ruleset

type HeroTrait int

const (
	TraitUnknown HeroTrait = iota

	TraitStartingHealthBonus // +3 max hp
	TraitStratingEnergyBonus // +5 max mp
	TraitStartingHybridBonus // +2 max hp, +3 max mp
	TraitHealthBonus         // +1 max hp
	TraitEnergyBonus         // +1 max mp
)
