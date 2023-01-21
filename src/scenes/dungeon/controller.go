package dungeon

import (
	"github.com/quasilyte/dicewind/assets"
	"github.com/quasilyte/dicewind/src/battle"
	"github.com/quasilyte/dicewind/src/gameui"
	"github.com/quasilyte/dicewind/src/ruleset"
	"github.com/quasilyte/dicewind/src/session"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmath"
)

type Controller struct {
	state *session.State

	scene *ge.Scene

	level *ruleset.DungeonLevel
	party *battle.Party

	infoScroll *gameui.InfoScroll
	hovered    any
	roomNodes  []*roomNode

	tiles []*gameui.UnitTile
}

func NewController(state *session.State, level *ruleset.DungeonLevel, party *battle.Party) *Controller {
	return &Controller{
		state: state,
		level: level,
		party: party,
	}
}

func (c *Controller) Init(scene *ge.Scene) {
	c.scene = scene
	c.initUI()
	c.updateUnitTiles()
}

func (c *Controller) Update(delta float64) {
	c.handleInput()
}

func (c *Controller) initUI() {
	ctx := c.scene.Context()
	bg := c.scene.NewRepeatedSprite(assets.ImageCryphBg, ctx.WindowWidth, ctx.WindowHeight)
	bg.Centered = false
	c.scene.AddGraphicsBelow(bg, 1)

	c.infoScroll = gameui.NewInfoScroll(gmath.Vec{X: 1104, Y: 28})
	c.scene.AddObject(c.infoScroll)

	pickedRooms := c.pickRooms()

	c.tiles = make([]*gameui.UnitTile, 6)
	for index := uint8(0); index < 6; index++ {
		tilePos := ruleset.TilePos{Alliance: 0, Index: index}
		offset := gameui.CalcUnitTilePos(tilePos)
		n := gameui.NewUnitTile(offset, tilePos)
		c.scene.AddObject(n)
		c.tiles[index] = n
		if index < 3 {
			pos := gmath.Vec{X: offset.X, Y: 896}
			rn := newRoomNode(pos, pickedRooms[index])
			c.scene.AddObject(rn)
			c.roomNodes = append(c.roomNodes, rn)
		}
	}
}

func (c *Controller) updateUnitTiles() {
	for i, tile := range c.tiles {
		tile.SetUnit(c.party.Heroes[i])
	}
}

func (c *Controller) pickRooms() []*ruleset.Room {
	if len(c.level.RoomDeck) < 3 {
		panic("not enough rooms to pick from")
	}
	gmath.Shuffle(c.scene.Rand(), c.level.RoomDeck)
	picked := make([]*ruleset.Room, 3)
	copy(picked, c.level.RoomDeck[len(c.level.RoomDeck)-3:len(c.level.RoomDeck)])
	c.level.RoomDeck = c.level.RoomDeck[:len(c.level.RoomDeck)-3]
	for _, r := range picked {
		if !r.Info.SingleVisit {
			c.level.RoomDeck = append(c.level.RoomDeck, r)
		}
	}
	return picked
}

func (c *Controller) handleInput() {
	cursorPos := c.state.MainInput.CursorPos()
	for _, rn := range c.roomNodes {
		if !rn.Rect.Contains(cursorPos) {
			continue
		}
		if c.hovered == rn {
			continue
		}
		c.hovered = rn
		c.infoScroll.SetText(rn.GetScrollText())
	}
}
