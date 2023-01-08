package ruleset

type TilePos struct {
	Alliance uint8
	Index    uint8
}

func GlobalTileIndexOf(alliance, index uint8) int {
	return (TilePos{Alliance: alliance, Index: index}).GlobalIndex()
}

func (pos TilePos) GlobalIndex() int {
	index := int(pos.Index)
	if pos.Alliance == 1 {
		index += 6
	}
	return index
}

func (pos TilePos) IsBackRow() bool {
	return pos.Index >= 3
}

func (pos TilePos) Col() int {
	col := int(pos.Index)
	if pos.IsBackRow() {
		col -= 3
	}
	return col
}

func (pos TilePos) OtherAlliance() TilePos {
	result := pos
	if pos.Alliance == 0 {
		result.Alliance = 1
	} else {
		result.Alliance = 0
	}
	return result
}

func (pos TilePos) LeftPos() TilePos {
	return TilePos{Alliance: pos.Alliance, Index: leftIndexTable[pos.Index]}
}

func (pos TilePos) RightPos() TilePos {
	return TilePos{Alliance: pos.Alliance, Index: rightIndexTable[pos.Index]}
}

func (pos TilePos) FrontRow() TilePos {
	if pos.IsBackRow() {
		return pos.OtherRow()
	}
	return pos
}

func (pos TilePos) BackRow() TilePos {
	if pos.IsBackRow() {
		return pos
	}
	return pos.OtherRow()
}

func (pos TilePos) OtherRow() TilePos {
	if pos.IsBackRow() {
		return TilePos{Alliance: pos.Alliance, Index: pos.Index - 3}
	}
	return TilePos{Alliance: pos.Alliance, Index: pos.Index + 3}
}

var leftIndexTable = [6]uint8{
	0: 0,
	1: 0,
	2: 1,
	3: 3,
	4: 3,
	5: 4,
}

var rightIndexTable = [6]uint8{
	0: 1,
	1: 2,
	2: 2,
	3: 4,
	4: 5,
	5: 5,
}
