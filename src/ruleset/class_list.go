package ruleset

import (
	"fmt"

	assets "github.com/quasilyte/dicewind/assets"
)

func HeroClassByName(name string) *HeroClass {
	for _, c := range HeroClassList {
		if c.Name == name {
			return c
		}
	}
	panic(fmt.Sprintf("unknown hero class %s", name))
}

var HeroClassList = []*HeroClass{
	{
		Name:      "Warrior",
		CardImage: assets.ImageHeroWarriorCard,
		HP:        6,
		MP:        2,
		Masteries: []MasteryKind{
			MasterySword,
		},
	},

	{
		Name:      "Sorcerer",
		CardImage: assets.ImageHeroSorcererCard,
		HP:        3,
		MP:        5,
		Masteries: []MasteryKind{
			MasteryStaff,
		},
	},
}
