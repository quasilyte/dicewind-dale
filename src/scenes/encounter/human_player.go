package encounter

import (
	"github.com/quasilyte/dicewind/src/battle"
	"github.com/quasilyte/dicewind/src/controls"
	"github.com/quasilyte/dicewind/src/ruleset"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/gesignal"
	"github.com/quasilyte/ge/input"
	"github.com/quasilyte/ge/tuple"
	"github.com/quasilyte/gmath"
)

type humanPlayer struct {
	input *input.Handler

	calc  *battle.Calculator
	board *battle.Board

	tileSelections [2][6]*tileSelectionAuraNode
	skillSlots     []*skillSlotNode
	skillTargeting int
	skillSelection *selectionAuraNode

	active bool

	unit    *battle.Unit
	actions []ruleset.Action

	EventActionsReady gesignal.Event[tuple.Value2[*battle.Unit, []ruleset.Action]]
}

func newHumanPlayer(h *input.Handler, calc *battle.Calculator, board *battle.Board) *humanPlayer {
	return &humanPlayer{
		input:          h,
		calc:           calc,
		board:          board,
		skillTargeting: -1,
	}
}

func (p *humanPlayer) Init(scene *ge.Scene) {
	for alliance := 0; alliance < 2; alliance++ {
		for pos := battle.TilePos(0); pos < 6; pos++ {
			a := newTileSelectionAuraNode(p.board.Tiles[alliance][pos].Pos)
			p.tileSelections[alliance][pos] = a
			scene.AddObjectBelow(a, 1)
		}
	}

	p.skillSelection = newSelectionAuraNode()
	scene.AddObjectAbove(p.skillSelection, 1)

	p.skillSlots = make([]*skillSlotNode, 4)
	pos := gmath.Vec{X: 1232, Y: 600}
	for i := range p.skillSlots {
		slot := newSkillSlotNode(pos)
		p.skillSlots[i] = slot
		scene.AddObject(slot)
		pos = pos.Add(gmath.Vec{X: 160 + 8})
	}
}

func (p *humanPlayer) IsDisposed() bool { return false }

func (p *humanPlayer) handleInput() {
	info, ok := p.input.JustPressedActionInfo(controls.ActionConfirm)
	if !ok || !info.HasPos() {
		return
	}

	cursorPos := info.Pos

	for i, slot := range p.skillSlots {
		if slot.Skill == nil || slot.Disabled {
			continue
		}
		if !slot.Rect.Contains(cursorPos) {
			continue
		}
		if p.skillTargeting == i {
			p.skillTargeting = -1
			p.updateTiles()
			continue
		}
		p.skillTargeting = i
		p.updateTiles()
		return
	}

	for alliance := 0; alliance < 2; alliance++ {
		for pos := battle.TilePos(0); pos < 6; pos++ {
			tile := p.tileSelections[alliance][pos]
			if tile.Action == ruleset.ActionNone {
				continue
			}
			if !tile.Rect.Contains(cursorPos) {
				continue
			}

			a := ruleset.Action{Kind: tile.Action}
			switch tile.Action {
			case ruleset.ActionMove, ruleset.ActionAttack:
				a.Value = int(pos)
			case ruleset.ActionSkill:
				a.Value = int(pos)
				a.SubKind = p.skillTargeting
			case ruleset.ActionGuard:
				// Do nothing.
			}
			p.actions = append(p.actions, a)
			p.sendActions()
			return
		}
	}
}

func (p *humanPlayer) Update(delta float64) {
	if !p.active {
		return
	}
	p.handleInput()
}

func (p *humanPlayer) sendActions() {
	p.active = false

	p.skillSelection.SetVisibility(false)

	for alliance := 0; alliance < 2; alliance++ {
		for pos := battle.TilePos(0); pos < 6; pos++ {
			tile := p.tileSelections[alliance][pos]
			tile.SetVisibility(false)
		}
	}
	for _, slot := range p.skillSlots {
		slot.SetSkill(nil)
	}

	p.EventActionsReady.Emit(tuple.New2(p.unit, p.actions))
}

func (p *humanPlayer) StartTurn(u *battle.Unit) {
	p.active = true
	p.unit = u
	p.actions = p.actions[:0]
	p.skillTargeting = -1

	skills := u.Skills()
	for i, skill := range skills {
		slot := p.skillSlots[i]
		slot.SetSkill(skill)
		slot.SetDisabled(skill.HealthCost >= u.HP || u.MP < skill.EnergyCost)
	}

	p.updateTiles()
}

func (p *humanPlayer) updateTiles() {
	if p.skillTargeting == -1 {
		p.skillSelection.SetVisibility(false)
	} else {
		p.skillSelection.SetVisibility(true)
		p.skillSelection.Pos = p.skillSlots[p.skillTargeting].pos
	}

	for alliance := 0; alliance < 2; alliance++ {
		for pos := battle.TilePos(0); pos < 6; pos++ {
			tile := p.tileSelections[alliance][pos]
			tile.SetVisibility(false)
			tile.Action = ruleset.ActionNone
			if p.skillTargeting == -1 {
				p.updateNormalTile(tile, alliance, pos)
			} else {
				p.updateSkillTile(tile, alliance, pos)
			}
		}
	}
}

func (p *humanPlayer) updateSkillTile(tile *tileSelectionAuraNode, alliance int, pos battle.TilePos) {
	skill := p.skillSlots[p.skillTargeting].Skill
	ok := (skill.CanTargetAlliedTile() && p.unit.Alliance == alliance) ||
		(skill.CanTargetEnemyTile() && p.unit.Alliance != alliance) ||
		(skill.TargetKind == ruleset.TargetEmptyAllied && p.unit.Alliance == alliance)
	if !ok {
		return
	}
	if p.calc.CanCastThere(p.unit, skill, p.unit.TilePos, pos) {
		tile.SetVisibility(true)
		tile.SetAction(ruleset.ActionSkill)
	}
}

func (p *humanPlayer) updateNormalTile(tile *tileSelectionAuraNode, alliance int, pos battle.TilePos) {
	u := p.unit
	if alliance == u.Alliance {
		if p.board.Tiles[alliance][pos].Unit == nil {
			tile.SetVisibility(true)
			tile.SetAction(ruleset.ActionMove)
		} else if pos == u.TilePos {
			tile.SetVisibility(true)
			tile.SetAction(ruleset.ActionGuard)
		}
		return
	}
	if p.board.Tiles[alliance][pos].Unit != nil && p.calc.CanAttackPos(u, u.TilePos, pos) {
		tile.SetVisibility(true)
		tile.SetAction(ruleset.ActionAttack)
	}
}
