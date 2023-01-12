package ruleset

import (
	"fmt"

	assets "github.com/quasilyte/dicewind/assets"
	"github.com/quasilyte/ge/resource"
)

// attack effects:
// - deal damage (can add magic damage to physical, etc)
// - ignore defense
// - hit multiple targets
// - add poison
// - add bleeding
// - burn mana
// - double strike
// - double roll (take best)
//
// debuffs:
// - stun
// - silence
// - dryad roots
// - damage reduction
//
// positive effects:
// - heal
//
// buffs:
// - damage boost
// - damage reduction
//
// random ideas:
// - collect charges before next attack
// - vampiric effect
// - attack and retreat at the same turn
// - move to a tile + attack from that tile
// - attack a random target with higher damage
// - swap pos with an ally
// - attack enemy and move it somewhere (pull towards you?)
//
// other skills:
// - summon creature
// - banish undead
// - unsummon creature
//

type Effect struct {
	Kind     EffectKind
	Source   EffectSource
	Value    any
	Duration int
}

type EffectSource int

const (
	SourceNone EffectSource = iota
	SourceMagical
	SourcePhysical
	SourcePsychological
)

type EffectKind int

const (
	EffectNone EffectKind = iota
	EffectDamage
	EffectDamageRow
	EffectAttack
	EffectRollBonus
	EffectPoison
	EffectManaBurn
	EffectWeakness
	EffectStun
	EffectSummonSkeleton
	EffectPoisonToHealth
)

type TargetKind int

const (
	TargetNone TargetKind = iota
	TargetSelf
	TargetAttackCandidate
	TargetEnemyRow
	TargetEnemyAny
	TargetEnemyMelee
	TargetEnemySpear
	TargetEmptyAllied
)

type Skill struct {
	Name            string
	Icon            resource.ImageID
	ImpactAnimation resource.ImageID
	CastSound       resource.AudioID
	EnergyCost      int
	HealthCost      int
	TargetKind      TargetKind
	TargetEffects   []Effect
}

func (s *Skill) CanTargetEnemyTile() bool {
	switch s.TargetKind {
	case TargetEnemyAny, TargetEnemyMelee, TargetEnemySpear, TargetAttackCandidate, TargetEnemyRow:
		return true
	default:
		return false
	}
}

func (s *Skill) CanTargetAlliedTile() bool {
	switch s.TargetKind {
	case TargetEmptyAllied, TargetSelf:
		return true
	default:
		return false
	}
}

func SkillByName(name string) *Skill {
	for _, s := range SkillList {
		if s.Name == name {
			return s
		}
	}
	panic(fmt.Sprintf("unknown skill %s", name))
}

var SkillList = []*Skill{
	{
		Name:            "True Strike",
		Icon:            assets.ImageSkillTrueStrike,
		ImpactAnimation: assets.ImageTrueStrike,
		CastSound:       assets.AudioScimitarAttack,
		EnergyCost:      1,
		TargetKind:      TargetAttackCandidate,
		TargetEffects: []Effect{
			{
				Kind:  EffectRollBonus,
				Value: 2,
			},
			{
				Kind:   EffectAttack,
				Source: SourcePhysical,
			},
		},
	},

	{
		Name:            "Consume Poison",
		Icon:            assets.ImageSkillConsumePoison,
		ImpactAnimation: assets.ImagePoisonExplosionReversed,
		CastSound:       assets.AudioPoisonExplosion,
		EnergyCost:      2,
		TargetKind:      TargetSelf,
		TargetEffects: []Effect{
			{
				Kind:   EffectPoisonToHealth,
				Source: SourcePhysical,
				Value:  2,
			},
		},
	},

	{
		Name:            "Summon Undead",
		Icon:            assets.ImageSkillSummonSkeleton,
		ImpactAnimation: assets.ImageDarkBoltExplosion,
		CastSound:       assets.AudioDarkExplosion,
		EnergyCost:      2,
		TargetKind:      TargetEmptyAllied,
		TargetEffects: []Effect{
			{
				Kind:   EffectSummonSkeleton,
				Source: SourceMagical,
			},
		},
	},

	{
		Name:            "Flame Strike",
		Icon:            assets.ImageSkillIconFlameStrike,
		ImpactAnimation: assets.ImageFlameStrike,
		CastSound:       assets.AudioFireExplosion,
		EnergyCost:      1,
		TargetKind:      TargetEnemyMelee,
		TargetEffects: []Effect{
			{
				Kind:   EffectDamage,
				Source: SourceMagical,
				Value:  DamageRange{1, 2, 2, 2, 2, 2},
			},
		},
	},

	{
		Name:            "Fireball",
		Icon:            assets.ImageSkillIconFireball,
		ImpactAnimation: assets.ImageFireExplosion,
		CastSound:       assets.AudioFireExplosion,
		EnergyCost:      3,
		TargetKind:      TargetEnemyAny,
		TargetEffects: []Effect{
			{
				Kind:   EffectDamage,
				Source: SourceMagical,
				Value:  DamageRange{1, 2, 3, 3, 3, 4},
			},
		},
	},

	{
		Name:            "Firestorm",
		Icon:            assets.ImageSkillIconFirestorm,
		ImpactAnimation: assets.ImageFirestormExplosion,
		CastSound:       assets.AudioFireExplosion,
		EnergyCost:      3,
		TargetKind:      TargetEnemyRow,
		TargetEffects: []Effect{
			{
				Kind:   EffectDamageRow,
				Source: SourceMagical,
				Value:  DamageRange{1, 2, 2, 2, 2, 3},
			},
		},
	},

	{
		Name:            "Hellfire",
		Icon:            assets.ImageSkillIconHellfire,
		ImpactAnimation: assets.ImageHellfireExplosion,
		CastSound:       assets.AudioFireExplosion,
		EnergyCost:      2,
		HealthCost:      1,
		TargetKind:      TargetEnemySpear,
		TargetEffects: []Effect{
			{
				Kind:   EffectDamage,
				Source: SourceMagical,
				Value:  DamageRange{2, 2, 2, 3, 4, 4},
			},
		},
	},
}
