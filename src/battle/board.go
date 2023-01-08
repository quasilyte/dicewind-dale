package battle

import "github.com/quasilyte/gmath"

type Board struct {
	// 3 4 5 ranged
	// 0 1 2 melee
	// -----
	// 0 1 2 melee
	// 3 4 5 ranged
	Tiles [2][6]Tile
}

type Tile struct {
	Unit *Unit

	Pos gmath.Vec
}

func NewBoard() *Board {
	b := &Board{}
	return b
}

func (b *Board) AddUnit(u *Unit, pos TilePos) {
	u.TilePos = pos
	b.Tiles[u.Alliance][pos].Unit = u
}

func (b *Board) WalkUnits(f func(u *Unit) bool) {
	for i := range b.Tiles {
		for _, tile := range b.Tiles[i] {
			if tile.Unit == nil {
				continue
			}
			if !f(tile.Unit) {
				return
			}
		}
	}
}
