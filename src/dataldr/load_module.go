package dataldr

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/quasilyte/dicewind/src/ruleset"
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

	rooms, err := loadModuleRooms(moduleRoot)
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
	filename := filepath.Join(moduleRoot, "module.json")
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var metadata moduleMetadata
	if err := json.Unmarshal(data, &metadata); err != nil {
		return nil, err
	}
	return &metadata, nil
}

type roomSchema struct {
	Name string `json:"name"`

	MaxVisits int `json:"max_visits"`

	Danger int `json:"danger"`

	Rewards map[string][]roomReward `json:"rewards"`
}

type roomReward struct {
	Kind string `json:"kind"`

	Value int `json:"value"`

	Scope string `json:"scope"`
}

func loadModuleRooms(moduleRoot string) ([]ruleset.RoomSchema, error) {
	filename := filepath.Join(moduleRoot, "rooms.json")
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
				switch raw.Kind {
				case "pick_skill":
					reward.Kind = ruleset.RewardPickSkill
				case "pick_potion":
					reward.Kind = ruleset.RewardPickPotion
				default:
					return rewards, fmt.Errorf("unknown %q reward kind", raw.Kind)
				}
				switch raw.Scope {
				case "every_party_member":
					reward.Scope = ruleset.RewardScopeEveryPartyMember
				case "shared":
					reward.Scope = ruleset.RewardScopeShared
				default:
					return rewards, fmt.Errorf("unknown %q reward scope", raw.Scope)
				}
				rewardList[i] = reward
			}

			rollRange := strings.Split(key, ",")
			for _, rollValue := range rollRange {
				v, err := strconv.Atoi(rollValue)
				if err != nil {
					return rewards, fmt.Errorf("bad roll range value: %q", rollValue)
				}
				if v < 1 || v > 6 {
					return rewards, fmt.Errorf("roll value %q is not in 1-6 range", rollValue)
				}
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
	convertRoom := func(raw roomSchema) (ruleset.RoomSchema, error) {
		room := ruleset.RoomSchema{
			Name:      raw.Name,
			MaxVisits: raw.MaxVisits,
			Danger:    raw.Danger,
		}
		rewards, err := convertRewards(raw.Rewards)
		if err != nil {
			return room, err
		}
		room.Rewards = rewards
		return room, nil
	}

	rooms := make([]ruleset.RoomSchema, len(rawRooms))
	for i, raw := range rawRooms {
		room, err := convertRoom(raw)
		if err != nil {
			return rooms, fmt.Errorf("invalid %q room data: %w", raw.Name, err)
		}
		rooms[i] = room
	}

	return rooms, nil
}
