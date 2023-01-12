package encounter

import (
	"github.com/quasilyte/ge/gesignal"
)

type effectGroup struct {
	numEffects int

	EventCompleted gesignal.Event[gesignal.Void]
}

func newEffectGroup() *effectGroup {
	return &effectGroup{}
}

func (g *effectGroup) AddEffect(e *effectNode) {
	g.numEffects++
	e.EventCompleted.Connect(g, func(gesignal.Void) {
		g.numEffects--
		if g.numEffects == 0 {
			g.EventCompleted.Emit(gesignal.Void{})
		}
	})
}

func (g *effectGroup) IsDisposed() bool {
	return false
}
