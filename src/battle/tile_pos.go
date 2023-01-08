package battle

type TilePos uint8

func (pos TilePos) IsBackRow() bool {
	return pos >= 3
}

func (pos TilePos) Col() int {
	col := int(pos)
	if pos.IsBackRow() {
		col -= 3
	}
	return col
}

func (pos TilePos) LeftPos() TilePos {
	return leftPosTable[pos]
}

func (pos TilePos) RightPos() TilePos {
	return rightPosTable[pos]
}

func (pos TilePos) OtherRow() TilePos {
	if pos.IsBackRow() {
		return pos - 3
	}
	return pos + 3
}

var leftPosTable = [6]TilePos{
	0: 0,
	1: 0,
	2: 1,
	3: 3,
	4: 3,
	5: 4,
}

var rightPosTable = [6]TilePos{
	0: 1,
	1: 2,
	2: 2,
	3: 4,
	4: 5,
	5: 5,
}
