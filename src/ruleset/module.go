package ruleset

type Module struct {
	Name string

	MaxDanger int

	Challenges []RoomChallenge

	Rooms []RoomSchema
}

type DungeonLevel struct {
	RoomsCleared int

	RoomDeck        []*Room
	VisitedRoomDeck []*Room
}

type Room struct {
	Module *Module

	Info RoomSchema
}

type RoomTrigger struct {
	Name string

	Challenges []*RoomChallenge
}

type RoomChallenge struct {
	Name string

	Kind RoomChallengeKind

	Value int
}

type RoomSchema struct {
	Name string

	SingleVisit bool

	Danger int

	Triggers [6]RoomTrigger

	Rewards [6][]RoomReward
}

type RoomReward struct {
	Kind RoomRewardKind

	Value int

	Scope RoomRewardScope
}

type RoomChallengeKind int

const (
	ChallengeUnknown RoomChallengeKind = iota
	ChallengeTrap
	ChallengePickableLock
)

type RoomRewardKind int

const (
	RewardUnknown RoomRewardKind = iota
	RewardPickSkill
	RewardPickPotion
	RewardPickLoot
	RewardRestoreHealth
	RewardRestoreEnergy
)

type RoomRewardScope int

const (
	RewardScopeUnknown RoomRewardScope = iota
	RewardScopeEveryPartyMember
	RewardScopeSelectedPartyMember
	RewardScopeShared
)
