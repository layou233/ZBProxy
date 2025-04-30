package route

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/layou233/zbproxy/v3/adapter"
	"github.com/layou233/zbproxy/v3/common/mcprotocol"
	"github.com/layou233/zbproxy/v3/common/set"
	"github.com/layou233/zbproxy/v3/config"
)

type RuleMinecraftPlayerName struct {
	sets      []set.StringSet
	lowerCase bool
	config    *config.Rule
}

var _ Rule = (*RuleMinecraftPlayerName)(nil)

func NewMinecraftPlayerNameRule(newConfig *config.Rule, listMap map[string]set.StringSet) (Rule, error) {
	var playerNameParameter config.RuleParameterListableString
	err := json.Unmarshal(newConfig.Parameter, &playerNameParameter)
	if err != nil {
		return nil, fmt.Errorf("bad player name parameter %v: %w", newConfig.Parameter, err)
	}
	sets := []set.StringSet{
		{}, // new set for individual names
	}
	for _, i := range playerNameParameter.Lists {
		if strings.HasPrefix(i, parameterListPrefix) {
			i = strings.TrimPrefix(i, parameterListPrefix)
			nameSet, found := listMap[i]
			if !found {
				return nil, fmt.Errorf("list [%v] is not found", i)
			}
			sets = append(sets, nameSet)
		} else {
			sets[0].Add(i)
		}
	}
	return &RuleMinecraftPlayerName{
		sets:      sets,
		lowerCase: playerNameParameter.LowerCase,
		config:    newConfig,
	}, nil
}

func (r *RuleMinecraftPlayerName) Config() *config.Rule {
	return r.config
}

func (r *RuleMinecraftPlayerName) Match(metadata *adapter.Metadata) (match bool) {
	if metadata.Minecraft != nil {
		name := metadata.Minecraft.PlayerName
		if r.lowerCase {
			name = strings.ToLower(name)
		}
		for _, nameSet := range r.sets {
			match = nameSet.Has(name)
			if match {
				break
			}
		}
	}
	if r.config.Invert {
		match = !match
	}
	return
}

type RuleMinecraftStatus struct {
	config *config.Rule
}

var _ Rule = (*RuleMinecraftStatus)(nil)

func NewMinecraftStatusRule(newConfig *config.Rule) (Rule, error) {
	return &RuleMinecraftStatus{
		config: newConfig,
	}, nil
}

func (r *RuleMinecraftStatus) Config() *config.Rule {
	return r.config
}

func (r *RuleMinecraftStatus) Match(metadata *adapter.Metadata) (match bool) {
	if metadata.Minecraft != nil {
		match = metadata.Minecraft.NextState == mcprotocol.NextStateStatus
	}
	if r.config.Invert {
		match = !match
	}
	return
}
