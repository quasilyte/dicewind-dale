package ruleset

type Module struct {
	Name string

	Rooms []RoomSchema
}

type DungeonLevel struct {
	RoomsCleared int

	RoomDeck []*Room
}

type Room struct {
	NumVisits int

	Info RoomSchema
}

type RoomSchema struct {
	Name string

	MaxVisits int

	Danger int

	Rewards [6][]RoomReward
}

type RoomReward struct {
	Kind RoomRewardKind

	Value int

	Scope RoomRewardScope
}

type RoomRewardKind int

const (
	RewardUnknown RoomRewardKind = iota
	RewardPickSkill
	RewardPickPotion
)

type RoomRewardScope int

const (
	RewardScopeUnknown RoomRewardScope = iota
	RewardScopeEveryPartyMember
	RewardScopeShared
)
