package assets

import (
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/resource"
)

func registerImages(ctx *ge.Context) {
	imageResources := map[resource.ImageID]resource.ImageInfo{
		ImageNoise: {Path: "image/noise.png"},

		ImageHeroSorcererCard: {Path: "image/unit/hero_sorcerer_card.png"},
		ImageHeroWarriorCard:  {Path: "image/unit/hero_warrior_card.png"},

		ImageSkeletonCard: {Path: "image/unit/skeleton_card.png"},

		ImageGreyMinionCard: {Path: "image/unit/minion_card.png"},

		ImageGreyMinionArcherCard: {Path: "image/unit/minion_archer_card.png"},

		ImageDarkspawnCard: {Path: "image/unit/darkspawn_card.png"},

		ImageLurkingTerrorCard: {Path: "image/unit/lurking_terror_card.png"},

		ImageBruteCard: {Path: "image/unit/brute_card.png"},

		ImageRoomBorder:        {Path: "image/room_border.png"},
		ImageUnitBorder:        {Path: "image/unit_border.png"},
		ImageUnitCardBg:        {Path: "image/unit_card_bg.png"},
		ImageSelectionAura:     {Path: "image/selection_aura.png"},
		ImageTileSelectionAura: {Path: "image/tile_selection_aura.png"},
		ImageSkillBorder:       {Path: "image/skill_border.png"},

		ImageSkillTrueStrike:      {Path: "image/skill/true_strike.png"},
		ImageSkillConsumePoison:   {Path: "image/skill/consume_poison.png"},
		ImageSkillSummonSkeleton:  {Path: "image/skill/summon_skeleton.png"},
		ImageSkillIconFirestorm:   {Path: "image/skill/firestorm.png"},
		ImageSkillIconFlameStrike: {Path: "image/skill/flame_strike.png"},
		ImageSkillIconFireball:    {Path: "image/skill/fireball.png"},
		ImageSkillIconHellfire:    {Path: "image/skill/hellfire.png"},

		ImageArrowAttack:    {Path: "image/arrow_attack.png", FrameWidth: 96},
		ImageBluntAttack:    {Path: "image/blunt_attack.png", FrameWidth: 96},
		ImageScimitarAttack: {Path: "image/scimitar_attack.png", FrameWidth: 96},
		ImageSwordAttack:    {Path: "image/sword_attack.png", FrameWidth: 96},
		ImageSpearAttack:    {Path: "image/spear_attack.png", FrameWidth: 96},
		ImageClawAttack:     {Path: "image/claw_attack.png", FrameWidth: 96},

		ImageCryphBg:     {Path: "image/dungeon/crypt_bg.png"},
		ImageCryptRoomBg: {Path: "image/dungeon/crypt_room_bg.png"},

		ImagePaperBg:     {Path: "image/paper_bg.png"},
		ImageEncounterBg: {Path: "image/encounter_bg.png"},
		ImagePoisonToken: {Path: "image/poison_token.png"},
		ImageHealthLevel: {Path: "image/health_level.png"},
		ImageEnergyLevel: {Path: "image/energy_level.png"},
		ImageHealthCost:  {Path: "image/health_cost.png"},
		ImageEnergyCost:  {Path: "image/energy_cost.png"},

		ImageTrueStrike:              {Path: "image/true_strike.png", FrameWidth: 96},
		ImageAcidSlingExplosion:      {Path: "image/acid_sling_explosion.png", FrameWidth: 156},
		ImagePoisonExplosion:         {Path: "image/poison_explosion.png", FrameWidth: 128},
		ImagePoisonExplosionReversed: {Path: "image/poison_explosion_reversed.png", FrameWidth: 128},
		ImageHellfireExplosion:       {Path: "image/hellfire_explosion.png", FrameWidth: 120},
		ImageFirestormExplosion:      {Path: "image/firestorm_explosion.png", FrameWidth: 128},
		ImageFireExplosion:           {Path: "image/fire_explosion.png", FrameWidth: 120},
		ImageDarkBoltExplosion:       {Path: "image/dark_bolt_explosion.png", FrameWidth: 64},
		ImageFlameStrike:             {Path: "image/flame_strike.png", FrameWidth: 128},
	}

	for id, res := range imageResources {
		ctx.Loader.ImageRegistry.Set(id, res)
		ctx.Loader.PreloadImage(id)
	}
}

const (
	ImageNone resource.ImageID = iota

	ImageNoise

	ImageHeroSorcererCard
	ImageHeroWarriorCard

	ImageSkeletonCard

	ImageGreyMinionArcherCard

	ImageGreyMinionCard

	ImageDarkspawnCard

	ImageLurkingTerrorCard

	ImageBruteCard

	ImageCryphBg
	ImageCryptRoomBg

	ImageRoomBorder
	ImageUnitBorder
	ImageUnitCardBg
	ImageSelectionAura
	ImageTileSelectionAura
	ImageSkillBorder

	ImageSkillTrueStrike
	ImageSkillConsumePoison
	ImageSkillSummonSkeleton
	ImageSkillIconFirestorm
	ImageSkillIconFlameStrike
	ImageSkillIconFireball
	ImageSkillIconHellfire

	ImageArrowAttack
	ImageBluntAttack
	ImageSwordAttack
	ImageSpearAttack
	ImageScimitarAttack
	ImageClawAttack

	ImagePaperBg
	ImageEncounterBg
	ImagePoisonToken
	ImageHealthLevel
	ImageEnergyLevel
	ImageHealthCost
	ImageEnergyCost

	ImageTrueStrike
	ImageAcidSlingExplosion
	ImagePoisonExplosion
	ImagePoisonExplosionReversed
	ImageFireExplosion
	ImageFirestormExplosion
	ImageHellfireExplosion
	ImageDarkBoltExplosion
	ImageFlameStrike
)
