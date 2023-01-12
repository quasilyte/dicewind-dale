package encounter

import (
	"github.com/quasilyte/dicewind/assets"
	"github.com/quasilyte/dicewind/src/battle"
	"github.com/quasilyte/dicewind/src/gameui"
	"github.com/quasilyte/dicewind/src/ruleset"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/gesignal"
)

type skillCastRunner struct {
	scene *ge.Scene

	event *battle.UnitSkillCastEvent
	board *battle.Board

	sourceLabel *gameui.Label

	delay float64
	fired bool

	EventCompleted gesignal.Event[gesignal.Void]
}

func newSkillCastRunner(board *battle.Board, e *battle.UnitSkillCastEvent) *skillCastRunner {
	return &skillCastRunner{
		event: e,
		board: board,
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
		if r.event.Skill.CastSound != assets.AudioNone {
			r.scene.Audio().PlaySound(r.event.Skill.CastSound)
		}
		if r.event.Skill.ImpactAnimation != assets.ImageNone {
			r.spawnEffects()
		} else {
			r.onAnimationCompleted(gesignal.Void{})
		}
	}
}

func (r *skillCastRunner) spawnEffects() {
	switch r.event.Skill.TargetKind {
	case ruleset.TargetEnemyRow:
		g := newEffectGroup()
		for col := 0; col < 3; col++ {
			targetPos := r.board.Tiles[r.event.Target.WithCol(col).GlobalIndex()].Pos
			e := newEffectNode(targetPos, r.event.Skill.ImpactAnimation)
			g.AddEffect(e)
			r.scene.AddObject(e)
		}
		g.EventCompleted.Connect(r, r.onAnimationCompleted)

	default:
		targetPos := r.board.Tiles[r.event.Target.GlobalIndex()].Pos
		e := newEffectNode(targetPos, r.event.Skill.ImpactAnimation)
		e.EventCompleted.Connect(r, r.onAnimationCompleted)
		r.scene.AddObject(e)
	}
}
