package dungeon

import "github.com/quasilyte/dicewind/src/ruleset"

func GenerateLevel(m *ruleset.Module) *ruleset.DungeonLevel {
	level := &ruleset.DungeonLevel{
		RoomDeck: make([]*ruleset.Room, 0, len(m.Rooms)),
	}

	for _, r := range m.Rooms {
		level.RoomDeck = append(level.RoomDeck, &ruleset.Room{
			Info: r,
		})
	}

	return level
}
