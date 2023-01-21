package main

import (
	"flag"
	"os"
	"strings"
	"time"

	"github.com/quasilyte/dicewind/assets"
	"github.com/quasilyte/dicewind/src/battle"
	"github.com/quasilyte/dicewind/src/controls"
	"github.com/quasilyte/dicewind/src/ruleset"
	"github.com/quasilyte/dicewind/src/scenes/encounter"
	"github.com/quasilyte/dicewind/src/session"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/langs"

	_ "image/png"
)

func main() {
	state := &session.State{}
	flag.StringVar(&state.AddonDir, "addon-dir", "$DICEWIND_DALE/_addon",
		"extra game data directory")
	flag.Parse()

	if strings.Contains(state.AddonDir, "$DICEWIND_DALE") {
		state.AddonDir = strings.ReplaceAll(state.AddonDir, "$DICEWIND_DALE", os.Getenv("DICEWIND_DALE"))
	}

	ctx := ge.NewContext()
	ctx.Rand.SetSeed(time.Now().Unix())
	ctx.GameName = "dicewind_dale"
	ctx.WindowTitle = "Dicewind Dale"
	ctx.WindowWidth = 1920
	ctx.WindowHeight = 1080
	ctx.FullScreen = true

	assets.Register(ctx)
	controls.BindKeymap(ctx, state)

	languages := ge.InferLanguages()
	preferredDict := assets.RawDictEn
	selectedLang := "en"
	for _, l := range languages {
		if l == "ru" {
			preferredDict = assets.RawDictRu
			selectedLang = "ru"
			break
		}
	}
	dict, err := langs.ParseDictionary(selectedLang, 2, ctx.Loader.LoadRaw(preferredDict))
	if err != nil {
		panic(err)
	}
	ctx.Dict = dict

	// m, err := dataldr.LoadModule(dict, filepath.Join(state.AddonDir, "crypt"))
	// if err != nil {
	// 	panic(err)
	// }
	// level := dungeon.GenerateLevel(m)
	// party := createTestParty()
	// if err := ge.RunGame(ctx, dungeon.NewController(state, level, party)); err != nil {
	// 	panic(err)
	// }

	if err := ge.RunGame(ctx, encounter.NewController(state)); err != nil {
		panic(err)
	}
}

func createTestParty() *battle.Party {
	warriorHero := &ruleset.Hero{
		Name:      "Alpha",
		CardImage: assets.ImageHeroWarriorCard,
		Traits: []ruleset.HeroTrait{
			ruleset.TraitStartingHealthBonus,
		},
		Weapon: &ruleset.HeroWeapon{
			Class: ruleset.WeaponByName("Sword"),
		},
		Armor: &ruleset.HeroArmor{
			Class: ruleset.ArmorByName("Mercenary Armor"),
		},
		CurrentSkills: []*ruleset.Skill{
			ruleset.SkillByName("True Strike"),
			ruleset.SkillByName("Consume Poison"),
		},
	}
	warriorHero.CurrentHP = warriorHero.MaxHP()
	warriorHero.CurrentMP = warriorHero.MaxMP()
	sorcHero := &ruleset.Hero{
		Name:      "Beta",
		CardImage: assets.ImageHeroSorcererCard,
		Traits: []ruleset.HeroTrait{
			ruleset.TraitStratingEnergyBonus,
		},
		Weapon: &ruleset.HeroWeapon{
			Class: ruleset.WeaponByName("Staff"),
		},
		CurrentSkills: []*ruleset.Skill{
			ruleset.SkillByName("Flame Strike"),
			ruleset.SkillByName("Fireball"),
			ruleset.SkillByName("Firestorm"),
			ruleset.SkillByName("Hellfire"),
		},
	}
	sorcHero.CurrentHP = sorcHero.MaxHP()
	sorcHero.CurrentMP = sorcHero.MaxMP()

	party := &battle.Party{}
	party.Heroes[0] = battle.NewHeroUnit(0, warriorHero)
	party.Heroes[4] = battle.NewHeroUnit(0, sorcHero)
	return party
}
