package ruleset

import (
	"fmt"

	assets "github.com/quasilyte/dicewind/assets"
)

var (
	aatk    = Action{Kind: ActionAttack}
	adef    = Action{Kind: ActionGuard}
	amove   = Action{Kind: ActionMove}
	askill1 = Action{Kind: ActionSkill, SubKind: 0}
	askill2 = Action{Kind: ActionSkill, SubKind: 1}
)

func MonsterByName(name string) *Monster {
	for _, m := range MonsterList {
		if m.Name == name {
			return m
		}
	}
	panic(fmt.Sprintf("unknown monster %s", name))
}

var MonsterList = []*Monster{
	{
		Name:      "Skeleton",
		CardImage: assets.ImageSkeletonCard,
		Level:     2,
		Reward:    3,
		HP:        3,
		Weapon:    WeaponSword,
		Damage:    DamageRange{0, 1, 1, 1, 1, 2},
		Reach:     ReachMelee,
		Action:    ActionRange{adef, adef, adef, aatk, aatk, aatk},
	},

	{
		Name:      "Grey Minion",
		CardImage: assets.ImageGreyMinionCard,
		Level:     1,
		Reward:    2,
		HP:        2,
		Weapon:    WeaponScimitar,
		Damage:    DamageRange{0, 0, 0, 1, 1, 1},
		Reach:     ReachMelee,
		Action:    ActionRange{amove, amove, adef, aatk, aatk, aatk},
		Race:      RaceGreyMinion,
	},

	{
		Name:      "Darkspawn",
		CardImage: assets.ImageDarkspawnCard,
		Level:     2,
		Reward:    4,
		HP:        3,
		Weapon:    WeaponClaws,
		Damage:    DamageRange{0, 1, 2, 2, 2, 2},
		Reach:     ReachMelee,
		Action:    ActionRange{adef, aatk, aatk, aatk, aatk, aatk},
	},

	{
		Name:      "Grey Minion Archer",
		CardImage: assets.ImageGreyMinionArcherCard,
		Level:     2,
		Reward:    3,
		HP:        2,
		Weapon:    WeaponBow,
		Damage:    DamageRange{0, 0, 1, 1, 1, 2},
		Reach:     ReachRanged,
		Action:    ActionRange{amove, aatk, aatk, aatk, askill1, askill1},
		Skills: []*Skill{
			{
				Name:            "Poison Arrow",
				ImpactAnimation: assets.ImagePoisonExplosion,
				CastSound:       assets.AudioPoisonExplosion,
				TargetKind:      TargetEnemyAny,
				TargetEffects: []Effect{
					{
						Kind:   EffectPoison,
						Source: SourcePhysical,
						Value:  2,
					},
				},
			},
		},
		Race: RaceGreyMinion,
	},

	{
		Name:   "Emissary of Despair",
		Level:  3,
		Reward: 5,
		HP:     5,
		Damage: DamageRange{0, 0, 1, 2, 3, 4},
		Reach:  ReachMelee,
		Action: ActionRange{adef, aatk, aatk, aatk, aatk, askill1},
		Skills: []*Skill{
			{
				Name: "Despair",
				TargetEffects: []Effect{
					{
						Kind:   EffectManaBurn,
						Source: SourcePsychological,
						Value:  2,
					},
				},
			},
		},
	},

	{
		Name:   "Grey Minion Veteran",
		Level:  3,
		Reward: 5,
		HP:     4,
		Damage: DamageRange{0, 1, 1, 2, 2, 2},
		Reach:  ReachMelee,
		Action: ActionRange{amove, aatk, aatk, aatk, aatk, aatk},
		Race:   RaceGreyMinion,
	},

	{
		Name:      "Lurking Terror",
		CardImage: assets.ImageLurkingTerrorCard,
		Level:     3,
		Reward:    6,
		HP:        4,
		Weapon:    WeaponClaws,
		Damage:    DamageRange{1, 1, 1, 2, 2, 3},
		Reach:     ReachMelee,
		Action:    ActionRange{aatk, aatk, aatk, askill1, askill1, askill1},
		Skills: []*Skill{
			{
				Name:            "Acid Sling",
				ImpactAnimation: assets.ImageAcidSlingExplosion,
				CastSound:       assets.AudioAcidSlingExplosion,
				EnergyCost:      2,
				TargetKind:      TargetEnemySpear,
				TargetEffects: []Effect{
					{
						Kind:   EffectDamage,
						Source: SourcePhysical,
						Value:  DamageRange{1, 2, 2, 2, 2, 2},
					},
				},
			},
		},
	},

	{
		Name:   "Grey Minion Warlock",
		Level:  4,
		Reward: 8,
		HP:     3,
		Damage: DamageRange{1, 1, 2, 2, 3, 3},
		Reach:  ReachRanged,
		Action: ActionRange{amove, amove, aatk, askill1, askill1, askill2},
		Race:   RaceGreyMinion,
		Skills: []*Skill{
			{
				Name: "Scorched Ground",
				TargetEffects: []Effect{
					{
						Kind:     EffectDamage,
						Source:   SourceMagical,
						Value:    2,
						Duration: 3,
					},
				},
			},
			{
				Name: "Weakness",
				TargetEffects: []Effect{
					{
						Kind:     EffectWeakness,
						Source:   SourceMagical,
						Duration: 3,
					},
				},
			},
		},
	},

	{
		Name:      "Brute",
		CardImage: assets.ImageBruteCard,
		Level:     4,
		Reward:    9,
		HP:        7,
		Weapon:    WeaponBlunt,
		Damage:    DamageRange{0, 0, 2, 2, 3, 3},
		Reach:     ReachMeleeFront,
		// Action:    ActionRange{aatk, aatk, aatk, aatk, askill1, askill1},
		Action: ActionRange{aatk, aatk, aatk, aatk, aatk, aatk},
		Skills: []*Skill{
			{
				Name: "Stun",
				TargetEffects: []Effect{
					{
						Kind:     EffectStun,
						Source:   SourcePhysical,
						Duration: 2,
					},
				},
			},
		},
	},
}
