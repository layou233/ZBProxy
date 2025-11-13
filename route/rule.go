package route

import (
	"errors"
	"fmt"
	"strings"

	"github.com/layou233/zbproxy/v3/adapter"
	"github.com/layou233/zbproxy/v3/common"
	"github.com/layou233/zbproxy/v3/common/set"
	"github.com/layou233/zbproxy/v3/config"

	"github.com/phuslu/log"
)

const (
	parameterListPrefix = "list:"
	typeCustomPrefix    = "custom:"
)

var ErrRuleTypeNotFound = errors.New("rule type is not found")

type Rule interface {
	Config() *config.Rule
	Match(metadata *adapter.Metadata) bool
}

func NewRule(logger *log.Logger, ruleConfig *config.Rule, listMap map[string]set.StringSet, ruleRegistry map[string]CustomRuleInitializer) (Rule, error) {
	switch ruleConfig.Type {
	case config.RuleTypeAlways:
		return &RuleAlways{ruleConfig}, nil
	case config.RuleTypeAnd:
		return NewLogicalAndRule(logger, ruleConfig, listMap, ruleRegistry)
	case config.RuleTypeOr:
		return NewLogicalOrRule(logger, ruleConfig, listMap, ruleRegistry)
	case config.RuleTypeServiceName:
		return NewServiceNameRule(ruleConfig, listMap)
	case config.RuleTypeSourceIPVersion:
		return NewSourceIPVersionRule(ruleConfig)
	case config.RuleTypeSourceIP:
		return NewSourceIPRule(ruleConfig, listMap)
	case config.RuleTypeSourcePort:
		return NewSourcePortRule(ruleConfig)
	case config.RuleTypeMinecraftHostname:
		return NewMinecraftHostnameRule(ruleConfig, listMap)
	case config.RuleTypeMinecraftPlayerName:
		return NewMinecraftPlayerNameRule(ruleConfig, listMap)
	case config.RuleTypeMinecraftStatus:
		return NewMinecraftStatusRule(ruleConfig)
	case config.RuleTypeMinecraftTransfer:
		return NewMinecraftTransferRule(ruleConfig)
	}
	if len(ruleRegistry) > 0 && strings.HasPrefix(ruleConfig.Type, typeCustomPrefix) {
		typeName := strings.TrimPrefix(ruleConfig.Type, typeCustomPrefix)
		ruleInitializer, found := ruleRegistry[typeName]
		if !found {
			return nil, fmt.Errorf("unknown custom rule type: %s", typeName)
		}
		return ruleInitializer(logger, ruleConfig, listMap)
	}
	return nil, common.Cause("type ["+ruleConfig.Type+"]: ", ErrRuleTypeNotFound)
}
