package ruleset

import "fmt"

func WeaponByName(name string) *HeroWeaponClass {
	for _, w := range WeaponList {
		if w.Name == name {
			return w
		}
	}
	panic(fmt.Sprintf("unknown weapon %s", name))
}

func ArmorByName(name string) *HeroArmorClass {
	for _, a := range ArmorList {
		if a.Name == name {
			return a
		}
	}
	panic(fmt.Sprintf("unknown armor %s", name))
}

var ArmorList = []*HeroArmorClass{
	{
		// Offers minimal protection for the wearer.
		// The defense increase effects are minimal,
		// but there are no negative side effects.
		Name:  "Mercenary Armor",
		Level: 1,
		Effects: []ItemEffect{
			{Kind: ItemEffectAttackerRollReduction, Value: 1},
		},
	},
}

var WeaponList = []*HeroWeaponClass{
	{
		Name:   "Sword",
		Level:  1,
		Kind:   WeaponSword,
		Damage: DamageRange{0, 1, 1, 1, 2, 3},
		Reach:  ReachMelee,
	},

	{
		Name:   "Scimitar",
		Level:  1,
		Kind:   WeaponScimitar,
		Damage: DamageRange{0, 0, 1, 2, 2, 3},
		Reach:  ReachMelee,
	},

	{
		Name:   "Staff",
		Level:  1,
		Kind:   WeaponBlunt,
		Damage: DamageRange{0, 0, 1, 1, 1, 2},
		Reach:  ReachMelee,
	},
}
