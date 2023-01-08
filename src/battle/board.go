package battle

import (
	"github.com/quasilyte/dicewind/src/ruleset"
	"github.com/quasilyte/gmath"
)

type Board struct {
	// 3 4 5 ranged
	// 0 1 2 melee
	// -----
	// 0 1 2 melee
	// 3 4 5 ranged
	Tiles [2 * 6]Tile
}

type Tile struct {
	Unit *Unit

	Pos gmath.Vec

	TilePos ruleset.TilePos
}

func NewBoard() *Board {
	b := &Board{}
	for alliance := uint8(0); alliance < 2; alliance++ {
		for index := uint8(0); index < 6; index++ {
			pos := ruleset.TilePos{Alliance: alliance, Index: index}
			b.Tiles[pos.GlobalIndex()].TilePos = pos
		}
	}
	return b
}

func (b *Board) GetTile(alliance, index uint8) *Tile {
	return &b.Tiles[ruleset.GlobalTileIndexOf(alliance, index)]
}

func (b *Board) AddUnit(u *Unit, pos ruleset.TilePos) {
	u.TilePos = pos
	b.Tiles[pos.GlobalIndex()].Unit = u
}

func (b *Board) WalkTiles(f func(t *Tile) bool) {
	for i := range b.Tiles {
		tile := &b.Tiles[i]
		if !f(tile) {
			return
		}
	}
}

func (b *Board) WalkTeamUnits(alliance uint8, f func(u *Unit) bool) {
	offset := 0
	if alliance == 1 {
		offset += 6
	}
	for i := 0; i < 6; i++ {
		tile := &b.Tiles[i+offset]
		if tile.Unit == nil {
			continue
		}
		if !f(tile.Unit) {
			return
		}
	}
}

func (b *Board) WalkUnits(f func(u *Unit) bool) {
	for i := range b.Tiles {
		tile := &b.Tiles[i]
		if tile.Unit == nil {
			continue
		}
		if !f(tile.Unit) {
			return
		}
	}
}
