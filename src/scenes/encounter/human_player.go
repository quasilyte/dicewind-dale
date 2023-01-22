package encounter

import (
	"strings"

	"github.com/quasilyte/dicewind/src/battle"
	"github.com/quasilyte/dicewind/src/controls"
	"github.com/quasilyte/dicewind/src/gameui"
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
	nodes *boardNodes

	skillSlots     []*skillSlotNode
	skillTargeting int
	skillSelection *selectionAuraNode

	active bool

	unit    *battle.Unit
	actions []ruleset.Action

	EventActionsReady gesignal.Event[tuple.Value2[*battle.Unit, []ruleset.Action]]
}

func newHumanPlayer(h *input.Handler, calc *battle.Calculator, board *battle.Board, nodes *boardNodes) *humanPlayer {
	return &humanPlayer{
		input:          h,
		calc:           calc,
		board:          board,
		nodes:          nodes,
		skillTargeting: -1,
	}
}

func (p *humanPlayer) Init(scene *ge.Scene) {
	p.skillSelection = newSelectionAuraNode()
	scene.AddObjectAbove(p.skillSelection, 1)

	keys := []string{
		p.input.ActionKeyNames(controls.ActionSkill1, input.KeyboardInput)[0],
		p.input.ActionKeyNames(controls.ActionSkill2, input.KeyboardInput)[0],
		p.input.ActionKeyNames(controls.ActionSkill3, input.KeyboardInput)[0],
		p.input.ActionKeyNames(controls.ActionSkill4, input.KeyboardInput)[0],
	}

	p.skillSlots = make([]*skillSlotNode, 4)
	pos := gmath.Vec{X: 1232, Y: 600}
	for i := range p.skillSlots {
		slot := newSkillSlotNode(strings.ToUpper(keys[i]), pos)
		p.skillSlots[i] = slot
		scene.AddObject(slot)
		pos = pos.Add(gmath.Vec{X: 160 + 8})
	}
}

func (p *humanPlayer) IsDisposed() bool { return false }

func (p *humanPlayer) Update(delta float64) {
	if !p.active {
		return
	}
	p.handleInput()
}

func (p *humanPlayer) handleInput() {
	skillActions := [...]input.Action{
		controls.ActionSkill1,
		controls.ActionSkill2,
		controls.ActionSkill3,
		controls.ActionSkill4,
	}
	for i, a := range &skillActions {
		slot := p.skillSlots[i]
		if slot.Skill == nil || slot.Disabled {
			continue
		}
		if p.input.ActionIsJustPressed(a) {
			p.activateSkill(i)
			return
		}
	}

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
		p.activateSkill(i)
		return
	}

	for _, tile := range p.nodes.tiles {
		if tile.GetAction() == ruleset.ActionNone {
			continue
		}
		if !tile.ContainsPos(cursorPos) {
			continue
		}

		a := ruleset.Action{Kind: tile.GetAction()}
		switch tile.GetAction() {
		case ruleset.ActionMove, ruleset.ActionAttack:
			a.Pos = tile.GetTilePos()
		case ruleset.ActionSkill:
			a.Pos = tile.GetTilePos()
			a.SubKind = p.skillTargeting
		case ruleset.ActionGuard:
			// Do nothing.
		}
		p.actions = append(p.actions, a)
		p.sendActions()
		return
	}
}

func (p *humanPlayer) activateSkill(i int) {
	if p.skillTargeting == i {
		p.skillTargeting = -1
		p.updateTiles()
		return
	}
	p.skillTargeting = i
	p.updateTiles()
}

func (p *humanPlayer) sendActions() {
	p.active = false

	p.skillSelection.SetVisibility(false)

	for _, tile := range p.nodes.tiles {
		tile.SetAction(ruleset.ActionNone)
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

	for _, tile := range p.nodes.tiles {
		tile.SetAction(ruleset.ActionNone)
		if p.skillTargeting == -1 {
			p.updateNormalTile(tile)
		} else {
			p.updateSkillTile(tile)
		}
	}
}

func (p *humanPlayer) updateSkillTile(tile *gameui.UnitTile) {
	pos := tile.GetTilePos()
	skill := p.skillSlots[p.skillTargeting].Skill
	ok := (skill.CanTargetAlliedTile() && p.unit.Alliance == pos.Alliance) ||
		(skill.CanTargetEnemyTile() && p.unit.Alliance != pos.Alliance) ||
		(skill.TargetKind == ruleset.TargetEmptyAllied && p.unit.Alliance == pos.Alliance)
	if !ok {
		return
	}
	if p.calc.CanCastThere(p.unit, skill, p.unit.TilePos, pos) {
		tile.SetAction(ruleset.ActionSkill)
	}
}

func (p *humanPlayer) updateNormalTile(tile *gameui.UnitTile) {
	pos := tile.GetTilePos()
	u := p.unit
	if pos.Alliance == u.Alliance {
		if p.board.Tiles[pos.GlobalIndex()].Unit == nil {
			tile.SetAction(ruleset.ActionMove)
		} else if pos == u.TilePos {
			tile.SetAction(ruleset.ActionGuard)
		}
		return
	}
	if p.board.Tiles[pos.GlobalIndex()].Unit != nil && p.calc.CanAttackPos(u, u.TilePos, pos) {
		tile.SetAction(ruleset.ActionAttack)
	}
}
