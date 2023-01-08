package encounter

import (
	"fmt"

	"github.com/quasilyte/dicewind/src/battle"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/gesignal"
)

type eventsRunner struct {
	scene *ge.Scene

	nodes *boardNodes

	active bool

	eventIndex int
	events     []battle.Event

	EventCompleted gesignal.Event[gesignal.Void]
}

func newEventsRunner(nodes *boardNodes) *eventsRunner {
	return &eventsRunner{nodes: nodes}
}

func (r *eventsRunner) Init(scene *ge.Scene) {
	r.scene = scene
}

func (r *eventsRunner) IsDisposed() bool { return false }

func (r *eventsRunner) Update(delta float64) {
	if !r.active {
		return
	}

	if r.eventIndex >= len(r.events) {
		r.active = false
		r.EventCompleted.Emit(gesignal.Void{})
		return
	}

	r.active = false
	e := r.events[r.eventIndex]
	r.eventIndex++
	r.runEvent(e)
}

func (r *eventsRunner) RunEvents(events []battle.Event) {
	r.active = true
	r.eventIndex = 0
	r.events = events
}

func (r *eventsRunner) runEvent(e battle.Event) {
	switch e := e.(type) {
	case *battle.UnitDamagedEvent:
		pos := ge.Pos{Base: &e.Unit.Pos}
		c := normalDamageColor
		if e.IsPoison {
			c = poisonDamageColor
		}
		r.scene.AddObject(newDamageScoreNode(e.Damage, pos.WithOffset(0, -32), c))
		r.active = true

	case *battle.UnitSkillCastEvent:
		targetNode := r.nodes.tiles[e.Target.GlobalIndex()]
		r2 := newSkillCastRunner(e, targetNode.body.Pos)
		r2.EventCompleted.Connect(r, r.subEventCompleted)
		r.scene.AddObject(r2)

	case *battle.UnitAttackEvent:
		r2 := newAttackRunner(e)
		r2.EventCompleted.Connect(r, r.subEventCompleted)
		r.scene.AddObject(r2)

	case *battle.UnitGuardEvent:
		r2 := newGuardRunner(e)
		r2.EventCompleted.Connect(r, r.subEventCompleted)
		r.scene.AddObject(r2)

	case *battle.UnitDefeatedEvent:
		n := r.nodes.tiles[e.Unit.TilePos.GlobalIndex()]
		gesignal.ConnectOneShot(&n.EventAnimationCompleted, nil, r.subEventCompleted)
		n.PlayUnitDefeat()

	case *battle.UnitMoveEvent:
		srcNode := r.nodes.tiles[e.From.GlobalIndex()]
		dstNode := r.nodes.tiles[e.To.GlobalIndex()]
		r2 := newMoveRunner(e, srcNode.body.Pos, dstNode.body.Pos, e.Unit.CardImage())
		r2.EventCompleted.Connect(r, r.subEventCompleted)
		r.scene.AddObject(r2)
		srcNode.SetUnit(nil)

	default:
		fmt.Printf("skip %T\n", e)
		r.active = true
	}
}

func (r *eventsRunner) subEventCompleted(gesignal.Void) {
	r.active = true
}
