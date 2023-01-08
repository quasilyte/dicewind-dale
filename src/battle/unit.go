package battle

import (
	assets "github.com/quasilyte/dicewind/assets"
	"github.com/quasilyte/dicewind/src/ruleset"
	"github.com/quasilyte/ge/resource"
	"github.com/quasilyte/gmath"
)

type Unit struct {
	Alliance int

	HP     int
	MP     int
	Poison int

	HumanControlled bool
	Guarding        bool

	TilePos TilePos
	Pos     gmath.Vec

	Monster *ruleset.Monster
	Hero    *ruleset.Hero
}

func NewHeroUnit(alliance int, h *ruleset.Hero) *Unit {
	return &Unit{
		Alliance:        alliance,
		Hero:            h,
		HP:              h.CurrentHP,
		MP:              h.CurrentMP,
		HumanControlled: true,
	}
}

func NewMonsterUnit(alliance int, m *ruleset.Monster) *Unit {
	return &Unit{
		Alliance:        alliance,
		Monster:         m,
		HP:              m.HP,
		HumanControlled: false,
	}
}

func (u *Unit) EnemyAlliance() int {
	if u.Alliance == 0 {
		return 1
	}
	return 0
}

func (u *Unit) IsMonster() bool {
	return u.Monster != nil
}

func (u *Unit) Skills() []*ruleset.Skill {
	if u.Monster != nil {
		return u.Monster.Skills
	}
	return u.Hero.CurrentSkills
}

func (u *Unit) WeaponKind() ruleset.WeaponKind {
	if u.Monster != nil {
		return u.Monster.Weapon
	}
	return u.Hero.WeaponKind()
}

func (u *Unit) AttackDamage() ruleset.DamageRange {
	if u.Monster != nil {
		return u.Monster.Damage
	}
	return u.Hero.DamageRange()
}

func (u *Unit) CardImage() resource.ImageID {
	if u.Monster != nil {
		return u.Monster.CardImage
	}
	return u.Hero.Class.CardImage
}

func (u *Unit) Name() string {
	if u.Monster != nil {
		return u.Monster.Name
	}
	return u.Hero.Name
}

func (u *Unit) WeaponMastery() ruleset.MasteryKind {
	if u.Monster != nil {
		return ruleset.MasteryNone
	}
	return u.Hero.Weapon.Class.Mastery
}

func (u *Unit) Masteries() []ruleset.MasteryKind {
	if u.Monster != nil {
		return nil
	}
	return u.Hero.Class.Masteries
}

func (u *Unit) Skill(i int) *ruleset.Skill {
	if u.Monster != nil {
		return u.Monster.Skills[i]
	}
	return u.Hero.CurrentSkills[i]
}

func (u *Unit) AttackSound() resource.AudioID {
	weaponKind := u.WeaponKind()
	switch weaponKind {
	case ruleset.WeaponClaws:
		return assets.AudioClawAttack
	case ruleset.WeaponScimitar:
		return assets.AudioScimitarAttack
	case ruleset.WeaponBlunt:
		return assets.AudioBluntAttack
	case ruleset.WeaponBow:
		return assets.AudioBowAttack
	case ruleset.WeaponSword:
		return assets.AudioSwordAttack
	default:
		panic("unknown weapon type")
	}
}

func (u *Unit) AttackImage() resource.ImageID {
	weaponKind := u.WeaponKind()
	switch weaponKind {
	case ruleset.WeaponClaws:
		return assets.ImageClawAttack
	case ruleset.WeaponScimitar:
		return assets.ImageScimitarAttack
	case ruleset.WeaponBlunt:
		return assets.ImageBluntAttack
	case ruleset.WeaponBow:
		return assets.ImageArrowAttack
	case ruleset.WeaponSword:
		return assets.ImageSwordAttack
	default:
		panic("unknown weapon type")
	}
}
