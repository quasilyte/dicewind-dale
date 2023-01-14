package dataldr

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/quasilyte/dicewind/src/jsonconfig"
	"github.com/quasilyte/dicewind/src/ruleset"
	"github.com/quasilyte/ge/xslices"
)

func LoadModule(moduleRoot string) (*ruleset.Module, error) {
	if _, err := os.Stat(moduleRoot); err != nil {
		return nil, err
	}

	var module ruleset.Module

	metadata, err := loadModuleMetadata(moduleRoot)
	if err != nil {
		return nil, err
	}
	module.Name = metadata.Name

	challenges, err := loadModuleChallenges(moduleRoot)
	if err != nil {
		return nil, err
	}
	module.Challenges = challenges

	rooms, err := loadModuleRooms(moduleRoot, challenges)
	if err != nil {
		return nil, err
	}
	module.Rooms = rooms

	return &module, nil
}

type moduleMetadata struct {
	Name string `json:"name"`
}

func loadModuleMetadata(moduleRoot string) (*moduleMetadata, error) {
	filename := filepath.Join(moduleRoot, "module.jsonc")
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var metadata moduleMetadata
	if err := jsonconfig.Unmarshal(data, &metadata); err != nil {
		return nil, err
	}
	return &metadata, nil
}

type roomTrigger struct {
	Name string `json:"name"`

	Challenges []string `json:"challenges"`
}

type roomChallenge struct {
	Name string `json:"name"`

	Kind string `json:"kind"`

	Value int `json:"value"`
}

func loadModuleChallenges(moduleRoot string) ([]ruleset.RoomChallenge, error) {
	filename := filepath.Join(moduleRoot, "challenges.jsonc")
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var rawChallenges []roomChallenge
	if err := jsonconfig.Unmarshal(data, &rawChallenges); err != nil {
		return nil, err
	}

	convertChallenge := func(raw roomChallenge) (ruleset.RoomChallenge, error) {
		challenge := ruleset.RoomChallenge{
			Name:  raw.Name,
			Value: raw.Value,
		}
		switch raw.Kind {
		case "pickable_lock":
			challenge.Kind = ruleset.ChallengePickableLock
		case "trap":
			challenge.Kind = ruleset.ChallengeTrap
		default:
			return challenge, fmt.Errorf("unknown %q challenge kind", raw.Kind)
		}
		return challenge, nil
	}

	challenges := make([]ruleset.RoomChallenge, len(rawChallenges))
	for i, raw := range rawChallenges {
		room, err := convertChallenge(raw)
		if err != nil {
			return challenges, fmt.Errorf("invalid %q challenge data: %w", raw.Name, err)
		}
		challenges[i] = room
	}

	return challenges, nil
}

type roomSchema struct {
	Name string `json:"name"`

	SingleVisit bool `json:"single_visit"`

	Danger int `json:"danger"`

	Triggers map[string]roomTrigger

	Rewards map[string][]roomReward `json:"rewards"`
}

type roomReward struct {
	Kind string `json:"kind"`

	Value int `json:"value"`

	Scope string `json:"scope"`
}

var rewardKinds = map[string]ruleset.RoomRewardKind{
	"pick_skill":     ruleset.RewardPickSkill,
	"pick_potion":    ruleset.RewardPickPotion,
	"pick_loot":      ruleset.RewardPickLoot,
	"restore_health": ruleset.RewardRestoreHealth,
	"restore_energy": ruleset.RewardRestoreEnergy,
}

var scopeKinds = map[string]ruleset.RoomRewardScope{
	"every_party_member":    ruleset.RewardScopeEveryPartyMember,
	"shared":                ruleset.RewardScopeShared,
	"selected_party_member": ruleset.RewardScopeSelectedPartyMember,
}

func loadModuleRooms(moduleRoot string, challenges []ruleset.RoomChallenge) ([]ruleset.RoomSchema, error) {
	filename := filepath.Join(moduleRoot, "rooms.jsonc")
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var rawRooms []roomSchema
	if err := json.Unmarshal(data, &rawRooms); err != nil {
		return nil, err
	}

	convertRewards := func(rawRewards map[string][]roomReward) ([6][]ruleset.RoomReward, error) {
		var rewards [6][]ruleset.RoomReward
		for key, rawList := range rawRewards {
			rewardList := make([]ruleset.RoomReward, len(rawList))
			for i, raw := range rawList {
				reward := ruleset.RoomReward{
					Value: raw.Value,
				}
				kind, ok := rewardKinds[raw.Kind]
				if !ok {
					return rewards, fmt.Errorf("unknown %q reward kind", raw.Kind)
				}
				reward.Kind = kind
				scope, ok := scopeKinds[raw.Scope]
				if !ok {
					return rewards, fmt.Errorf("unknown %q reward scope", raw.Scope)
				}
				reward.Scope = scope
				rewardList[i] = reward
			}

			rollValues, err := parseRollRange(key)
			if err != nil {
				return rewards, err
			}
			for _, v := range rollValues {
				rewards[v-1] = rewardList
			}
		}

		for i := range rewards {
			if len(rewards[i]) == 0 {
				return rewards, fmt.Errorf("no rewards specified for roll of %d", i+1)
			}
		}

		return rewards, nil
	}

	convertTriggers := func(rawTriggers map[string]roomTrigger, challenges []ruleset.RoomChallenge) ([6]ruleset.RoomTrigger, error) {
		var triggers [6]ruleset.RoomTrigger
		for key, raw := range rawTriggers {
			rollValues, err := parseRollRange(key)
			if err != nil {
				return triggers, err
			}
			trigger := ruleset.RoomTrigger{
				Name:       raw.Name,
				Challenges: make([]*ruleset.RoomChallenge, 0, len(raw.Challenges)),
			}
			for _, challengeName := range raw.Challenges {
				i := xslices.IndexWhere(challenges, func(c ruleset.RoomChallenge) bool {
					return c.Name == challengeName
				})
				if i == -1 {
					return triggers, fmt.Errorf("referenced undefined challenge %q", challengeName)
				}
				trigger.Challenges = append(trigger.Challenges, &challenges[i])
			}
			for _, v := range rollValues {
				triggers[v-1] = trigger
			}
		}
		return triggers, nil
	}

	convertRoom := func(raw roomSchema, challenges []ruleset.RoomChallenge) (ruleset.RoomSchema, error) {
		room := ruleset.RoomSchema{
			Name:        raw.Name,
			SingleVisit: raw.SingleVisit,
			Danger:      raw.Danger,
		}

		rewards, err := convertRewards(raw.Rewards)
		if err != nil {
			return room, err
		}
		room.Rewards = rewards

		triggers, err := convertTriggers(raw.Triggers, challenges)
		if err != nil {
			return room, err
		}
		room.Triggers = triggers

		return room, nil
	}

	rooms := make([]ruleset.RoomSchema, len(rawRooms))
	for i, raw := range rawRooms {
		room, err := convertRoom(raw, challenges)
		if err != nil {
			return rooms, fmt.Errorf("invalid %q room data: %w", raw.Name, err)
		}
		rooms[i] = room
	}

	return rooms, nil
}

func parseRollRange(s string) ([]int, error) {
	rollRange := strings.Split(s, ",")
	values := make([]int, 0, 6)
	for _, rollValue := range rollRange {
		v, err := strconv.Atoi(rollValue)
		if err != nil {
			return nil, fmt.Errorf("bad roll range value: %q", rollValue)
		}
		if v < 1 || v > 6 {
			return nil, fmt.Errorf("roll value %q is not in 1-6 range", rollValue)
		}
		values = append(values, v)
	}
	return values, nil
}
