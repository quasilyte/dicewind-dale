package encounter

import (
	"github.com/quasilyte/dicewind/assets"
	"github.com/quasilyte/dicewind/src/battle"
	"github.com/quasilyte/dicewind/src/gameui"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/gesignal"
	"github.com/quasilyte/gmath"
)

type skillCastRunner struct {
	scene *ge.Scene

	event     *battle.UnitSkillCastEvent
	targetPos gmath.Vec

	sourceLabel *gameui.Label

	delay float64
	fired bool

	EventCompleted gesignal.Event[gesignal.Void]
}

func newSkillCastRunner(e *battle.UnitSkillCastEvent, targetPos gmath.Vec) *skillCastRunner {
	return &skillCastRunner{
		event:     e,
		targetPos: targetPos,
	}
}

func (r *skillCastRunner) Init(scene *ge.Scene) {
	r.scene = scene

	r.sourceLabel = gameui.NewLabel(gameui.LabelConfig{
		Width:  196,
		Height: 64,
		Pos:    ge.Pos{Base: &r.event.Caster.Pos},
		Font:   assets.FontTiny,
		Text:   r.event.Skill.Name,
	})
	scene.AddObject(r.sourceLabel)

	r.delay = 0.5
}

func (r *skillCastRunner) IsDisposed() bool { return r.sourceLabel.IsDisposed() }

func (r *skillCastRunner) Dispose() {
	r.sourceLabel.Dispose()
}

func (r *skillCastRunner) onAnimationCompleted(gesignal.Void) {
	// pos := ge.Pos{Base: &r.event.Defender.Pos}
	// r.scene.AddObject(newDamageScoreNode(r.event.Damage, pos.WithOffset(0, -32)))

	r.EventCompleted.Emit(gesignal.Void{})
	r.Dispose()
}

func (r *skillCastRunner) Update(delta float64) {
	if r.fired {
		return
	}
	r.delay -= delta
	if r.delay <= 0 {
		r.fired = true
		r.scene.Audio().PlaySound(r.event.Skill.CastSound)
		e := newEffectNode(r.targetPos, r.event.Skill.ImpactAnimation)
		e.EventCompleted.Connect(r, r.onAnimationCompleted)
		r.scene.AddObject(e)
	}
}
