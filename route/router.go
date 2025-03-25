package route

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"
	"sync"

	"github.com/layou233/zbproxy/v3/adapter"
	"github.com/layou233/zbproxy/v3/common"
	"github.com/layou233/zbproxy/v3/common/bufio"
	"github.com/layou233/zbproxy/v3/common/set"
	"github.com/layou233/zbproxy/v3/config"
	"github.com/layou233/zbproxy/v3/protocol"

	"github.com/phuslu/log"
)

type RouterOptions struct {
	Config          *config.Router
	OutboundMap     map[string]adapter.Outbound
	ListMap         map[string]set.StringSet
	RuleRegistry    map[string]CustomRuleInitializer
	SnifferRegistry map[string]protocol.SnifferFunc
}

type Router struct {
	access          sync.RWMutex
	ctx             context.Context
	logger          *log.Logger
	outboundMap     map[string]adapter.Outbound
	listMap         map[string]set.StringSet
	ruleRegistry    map[string]CustomRuleInitializer
	snifferRegistry map[string]protocol.SnifferFunc
	rules           []Rule
	defaultOutbound adapter.Outbound
	started         bool
}

var _ adapter.Router = (*Router)(nil)

func (r *Router) Initialize(ctx context.Context, logger *log.Logger, options RouterOptions) error {
	if r.started {
		return errors.New("already initialized")
	}
	rules := make([]Rule, 0, len(options.Config.Rules))
	for i, ruleConfig := range options.Config.Rules {
		rule, err := NewRule(logger, ruleConfig, options.ListMap, options.RuleRegistry)
		if err != nil {
			return fmt.Errorf("initialize rule [index=%d]: %w", i, err)
		}
		rules = append(rules, rule)
	}
	if options.Config.DefaultOutbound != "" {
		defaultOutbound, err := r.findOutboundByName(options.OutboundMap, options.Config.DefaultOutbound)
		if err != nil {
			return common.Cause("default outbound is not found: ", err)
		}
		r.defaultOutbound = defaultOutbound
	} else {
		r.defaultOutbound, _ = protocol.NewOutbound(r.logger, &config.Outbound{
			Name: "default",
		})
		r.defaultOutbound.PostInitialize(r, r) // this is dangerous since the router is not fully initialized yet
	}
	r.ctx = ctx
	r.logger = logger
	r.outboundMap = options.OutboundMap
	r.listMap = options.ListMap
	r.ruleRegistry = options.RuleRegistry
	r.snifferRegistry = options.SnifferRegistry
	r.rules = rules
	r.started = true
	return nil
}

func (r *Router) HandleConnection(conn net.Conn, metadata *adapter.Metadata) {
	r.access.RLock()
	cachedConn := bufio.NewCachedConn(conn)
	outbound := r.defaultOutbound
	for i, rule := range r.rules {
		match := rule.Match(metadata)
		if match {
			r.logger.Trace().Str("id", metadata.ConnectionID).Int("rule_index", i).Msg("Rule matched")
			ruleConfig := rule.Config()
			// handle sniff
			if len(ruleConfig.Sniff) > 0 {
				protocol.Sniff(r.logger, cachedConn, metadata, r.snifferRegistry, ruleConfig.Sniff...)
			}
			// handle rewrite
			if ruleConfig.Rewrite.TargetAddress != "" {
				metadata.DestinationHostname = ruleConfig.Rewrite.TargetAddress
			}
			if ruleConfig.Rewrite.TargetPort > 0 {
				metadata.DestinationPort = ruleConfig.Rewrite.TargetPort
			}
			if ruleConfig.Rewrite.Minecraft != nil {
				if metadata.Minecraft == nil {
					r.logger.Debug().Str("id", metadata.ConnectionID).Int("rule_index", i).Msg("No Minecraft metadata, skipped rewrite")
				} else {
					if ruleConfig.Rewrite.Minecraft.Hostname != "" {
						metadata.Minecraft.RewrittenDestination = ruleConfig.Rewrite.Minecraft.Hostname
					}
					if ruleConfig.Rewrite.Minecraft.Port > 0 {
						metadata.Minecraft.RewrittenPort = ruleConfig.Rewrite.Minecraft.Port
					}
				}
			}
			// handle outbound
			if ruleConfig.Outbound != "" {
				var err error
				outbound, err = r.FindOutboundByName(ruleConfig.Outbound)
				if err != nil {
					r.logger.Error().Str("id", metadata.ConnectionID).Int("rule_index", i).
						Err(err).Msg("Failed to find outbound")
					conn.Close()
					r.access.RUnlock()
					return
				}
				break
			}
		}
	}

	if injectOutbound, isInject := outbound.(adapter.InjectOutbound); isInject {
		r.access.RUnlock()
		err := injectOutbound.InjectConnection(r.ctx, cachedConn, metadata)
		var logger *log.Entry
		if err == nil {
			logger = r.logger.Info()
		} else {
			logger = r.logger.Warn()
		}
		logger = logger.Str("id", metadata.ConnectionID).Str("outbound", outbound.Name())
		if err != nil {
			logger = logger.Err(err)
		}
		logger.Msg("Handled outbound connection")
		cachedConn.Close()
		return
	} else if metadata.DestinationHostname != "" && metadata.DestinationPort > 0 {
		destinationConn, err := adapter.DialContextWithMetadata(outbound, r.ctx, "tcp",
			net.JoinHostPort(metadata.DestinationHostname, strconv.FormatUint(uint64(metadata.DestinationPort), 10)),
			metadata)
		if err != nil {
			r.logger.Warn().Str("id", metadata.ConnectionID).Str("outbound", outbound.Name()).
				Err(err).Msg("Failed to dial outbound connection")
		}
		r.access.RUnlock()
		err = bufio.CopyConn(destinationConn, cachedConn)
		var logger *log.Entry
		if err == nil {
			logger = r.logger.Info()
		} else {
			logger = r.logger.Warn()
		}
		logger = logger.Str("id", metadata.ConnectionID).Str("outbound", outbound.Name())
		if err != nil {
			logger = logger.Err(err)
		}
		logger.Msg("Handled outbound connection")
		cachedConn.Close()
		return
	}
	r.logger.Info().Str("id", metadata.ConnectionID).Msg("Closed leaked connection")
	cachedConn.Close()
	r.access.RUnlock()
}

func (r *Router) FindOutboundByName(name string) (adapter.Outbound, error) {
	return r.findOutboundByName(r.outboundMap, name)
}

func (r *Router) findOutboundByName(outboundMap map[string]adapter.Outbound, name string) (adapter.Outbound, error) {
	switch name {
	case "REJECT":
		return rejectOutbound{}, nil
	case "RESET":
		return resetOutbound{}, nil
	}
	if outboundMap == nil {
		return nil, errors.New("outbounds are not initialized")
	}
	outbound, ok := outboundMap[name]
	if !ok {
		return nil, fmt.Errorf("outbound not found [%s]", name)
	}
	return outbound, nil
}

func (r *Router) FindListsByTag(tags []string) ([]set.StringSet, error) {
	if r.listMap == nil {
		return nil, errors.New("lists are not initialized")
	}
	lists := make([]set.StringSet, 0, len(tags))
	for _, tag := range tags {
		list, ok := r.listMap[tag]
		if !ok {
			return nil, fmt.Errorf("list not found [%s]", tag)
		}
		lists = append(lists, list)
	}
	return lists, nil
}

func (r *Router) UpdateConfig(newOptions RouterOptions) error {
	r.access.Lock()
	defer r.access.Unlock()
	r.started = false
	return r.Initialize(r.ctx, r.logger, newOptions)
}
