package dungeon

import (
	"fmt"
	"strings"

	"github.com/quasilyte/dicewind/assets"
	"github.com/quasilyte/dicewind/src/gameui"
	"github.com/quasilyte/dicewind/src/ruleset"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmath"
)

type roomNode struct {
	pos gmath.Vec

	scene *ge.Scene

	nameLabel *gameui.LabelWithShadow

	room *ruleset.Room

	bg     *ge.Sprite
	border *ge.Sprite

	Rect gmath.Rect
}

func newRoomNode(pos gmath.Vec, room *ruleset.Room) *roomNode {
	return &roomNode{pos: pos, room: room}
}

func (n *roomNode) Init(scene *ge.Scene) {
	n.scene = scene

	n.bg = scene.NewSprite(assets.ImageCryptRoomBg)
	n.bg.Pos.Base = &n.pos
	scene.AddGraphics(n.bg)

	n.border = scene.NewSprite(assets.ImageRoomBorder)
	n.border.Pos.Base = &n.pos
	scene.AddGraphics(n.border)

	n.nameLabel = gameui.NewLabelWithShadow(gameui.LabelWithShadowConfig{
		Font:  assets.FontTiny,
		Pos:   ge.Pos{Base: &n.pos, Offset: gmath.Vec{X: (-320 / 2) + 16, Y: 80}},
		Width: 320 - 32,
	})
	n.nameLabel.Text = scene.Dict().Get(n.room.Module.Name, "rooms", n.room.Info.Name)
	scene.AddObject(n.nameLabel)

	n.Rect.Min = n.border.AnchorPos().Resolve()
	n.Rect.Max = n.Rect.Min.Add(gmath.Vec{X: n.border.ImageWidth(), Y: n.border.ImageHeight()})
}

func (n *roomNode) IsDisposed() bool { return false }

func (n *roomNode) Update(delta float64) {
}

func (n *roomNode) GetScrollText() (title, rtitle, content string) {
	r := n.room
	dict := n.scene.Dict()

	title = dict.Get(r.Module.Name, "rooms", r.Info.Name)
	rtitle = fmt.Sprintf("%s: %d/%d", dict.GetTitleCase("word.danger"), r.Info.Danger, r.Module.MaxDanger)
	var lines []string
	lines = append(lines, dict.Get(r.Module.Name, "rooms", r.Info.Name, "description"))
	return title, rtitle, strings.Join(lines, "\n")
}
