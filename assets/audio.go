package assets

import (
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/resource"
)

func registerAudio(ctx *ge.Context) {
	audioResources := map[resource.AudioID]resource.Audio{
		AudioClawAttack:     {Path: "audio/claw_attack.wav", Volume: -0.4},
		AudioSwordAttack:    {Path: "audio/sword_attack.wav", Volume: -0.6},
		AudioBluntAttack:    {Path: "audio/blunt_attack.wav", Volume: -0.4},
		AudioScimitarAttack: {Path: "audio/scimitar_attack.wav", Volume: -0.8},
		AudioBowAttack:      {Path: "audio/bow_attack.wav", Volume: -0.8},

		AudioDarkExplosion:      {Path: "audio/dark_explosion.wav", Volume: -0.4},
		AudioAcidSlingExplosion: {Path: "audio/acid_sling_explosion.wav", Volume: -0.2},
		AudioPoisonExplosion:    {Path: "audio/poison_explosion.wav", Volume: -0.4},
		AudioFireExplosion:      {Path: "audio/fire_explosion.wav", Volume: -0.1},
	}
	for id, res := range audioResources {
		ctx.Loader.AudioRegistry.Set(id, res)
		ctx.Loader.PreloadAudio(id)
	}
}

const (
	AudioNone resource.AudioID = iota

	AudioClawAttack
	AudioSwordAttack
	AudioBluntAttack
	AudioScimitarAttack
	AudioBowAttack

	AudioDarkExplosion
	AudioAcidSlingExplosion
	AudioPoisonExplosion
	AudioFireExplosion
)
