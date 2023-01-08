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

var WeaponList = []*HeroWeaponClass{
	{
		Name:    "Sword",
		Level:   1,
		Kind:    WeaponSword,
		Mastery: MasterySword,
		Damage:  DamageRange{0, 1, 2, 2, 2, 2},
		Reach:   ReachMelee,
	},

	{
		Name:    "Staff",
		Level:   1,
		Kind:    WeaponBlunt,
		Mastery: MasteryStaff,
		Damage:  DamageRange{0, 0, 1, 1, 1, 2},
		Reach:   ReachMelee,
	},
}
