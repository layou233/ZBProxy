package zbproxy

import (
	"context"
	"errors"
	"time"

	"github.com/layou233/zbproxy/v3/adapter"
	"github.com/layou233/zbproxy/v3/common"
	"github.com/layou233/zbproxy/v3/config"
	"github.com/layou233/zbproxy/v3/protocol"
	"github.com/layou233/zbproxy/v3/route"
	"github.com/layou233/zbproxy/v3/service"

	"github.com/phuslu/log"
)

type Options struct {
	Config          *config.Root
	LogWriter       log.Writer
	ConfigFilePath  string
	RuleRegistry    map[string]route.CustomRuleInitializer
	SnifferRegistry map[string]protocol.SnifferFunc
	DisableReload   bool
}

type Instance struct {
	ctx             context.Context
	logger          *log.Logger
	config          *config.Root
	router          *route.Router
	serviceMap      map[string]adapter.Service
	outboundMap     map[string]adapter.Outbound
	ruleRegistry    map[string]route.CustomRuleInitializer
	snifferRegistry map[string]protocol.SnifferFunc
}

func NewInstance(ctx context.Context, options Options) (*Instance, error) {
	instance := &Instance{
		ctx: ctx,
		logger: &log.Logger{
			TimeFormat: "15:04:05",
			Writer:     options.LogWriter,
		},
		config:          options.Config,
		ruleRegistry:    options.RuleRegistry,
		snifferRegistry: options.SnifferRegistry,
	}
	if options.LogWriter == nil {
		instance.logger.Writer = &log.ConsoleWriter{
			ColorOutput:    true,
			EndWithMessage: true,
		}
	}
	if options.Config == nil {
		if options.ConfigFilePath != "" {
			newConfig, err := config.LoadConfigFromFile(
				ctx, options.ConfigFilePath, !options.DisableReload, instance.logger)
			if err != nil {
				return nil, err
			}
			instance.config = newConfig
			if !options.DisableReload {
				// configure handler here
				// if you provide config directly from option
				// then configure it by yourself like this
				newConfig.SetUpdateHandler(instance.UpdateConfig)
			}
		} else {
			return nil, errors.New("no config provided for ZBProxy")
		}
	}
	instance.logger.Level = instance.config.Log.Level

	return instance, nil
}

func (i *Instance) Logger() *log.Logger {
	return i.logger
}

func (i *Instance) Router() *route.Router {
	return i.router
}

func (i *Instance) Start() error {
	var err error
	startTime := time.Now()

	// initialize outbounds
	outboundMap := make(map[string]adapter.Outbound, len(i.config.Outbounds))
	for _, outboundConfig := range i.config.Outbounds {
		var outbound adapter.Outbound
		outbound, err = protocol.NewOutbound(i.logger, outboundConfig)
		if err != nil {
			return common.Cause("initialize outbound ["+outboundConfig.Name+"]: ", err)
		}
		outboundMap[outbound.Name()] = outbound
	}
	i.outboundMap = outboundMap

	// initialize router
	i.router = &route.Router{}
	err = i.router.Initialize(i.ctx, i.logger, route.RouterOptions{
		Config:          &i.config.Router,
		OutboundMap:     outboundMap,
		ListMap:         i.config.Lists,
		RuleRegistry:    i.ruleRegistry,
		SnifferRegistry: i.snifferRegistry,
	})
	if err != nil {
		return common.Cause("initialize router: ", err)
	}
	for _, outbound := range outboundMap {
		err = outbound.PostInitialize(i.router, i.router)
		if err != nil {
			return common.Cause("post initialize outbound ["+outbound.Name()+"]: ", err)
		}
	}

	// initialize services
	i.serviceMap = make(map[string]adapter.Service, len(i.config.Services))
	for _, serviceConfig := range i.config.Services {
		newService := service.NewService(i.logger, serviceConfig)
		newService.UpdateRouter(i.router)
		err = newService.Start(i.ctx)
		if err != nil {
			return common.Cause("start service ["+serviceConfig.Name+"]: ", err)
		}
		i.serviceMap[serviceConfig.Name] = newService
	}

	i.logger.Info().Str("duration", time.Now().Sub(startTime).String()).Msg("ZBProxy started")
	return nil
}

func (i *Instance) Reload() bool {
	return i.config.Reload()
}

func (i *Instance) UpdateConfig() {
	// update outbounds
	newOutboundMap := make(map[string]adapter.Outbound, len(i.config.Outbounds))
	for _, outboundConfig := range i.config.Outbounds {
		if oldOutbound, ok := i.outboundMap[outboundConfig.Name]; ok {
			err := oldOutbound.Reload(adapter.OutboundReloadOptions{
				Router: i.router,
				Config: outboundConfig,
				Lists:  i.config.Lists,
			})
			if err != nil {
				i.logger.Error().Str("outbound", outboundConfig.Name).Err(err).Msg("Error when updating outbounds")
				return
			}
			newOutboundMap[outboundConfig.Name] = oldOutbound
		} else {
			newOutbound, err := protocol.NewOutbound(i.logger, outboundConfig)
			if err != nil {
				i.logger.Error().Str("outbound", outboundConfig.Name).Err(err).Msg("Error when initializing outbounds")
				return
			}
			newOutboundMap[outboundConfig.Name] = newOutbound
		}
	}

	// update router
	err := i.router.UpdateConfig(route.RouterOptions{
		Config:          &i.config.Router,
		OutboundMap:     newOutboundMap,
		ListMap:         i.config.Lists,
		RuleRegistry:    i.ruleRegistry,
		SnifferRegistry: i.snifferRegistry,
	})
	if err != nil {
		i.logger.Error().Err(err).Msg("Error when updating router")
		return
	}

	// update services
	newServiceMap := make(map[string]adapter.Service, len(i.config.Services))
	for _, serviceConfig := range i.config.Services {
		if oldService, ok := i.serviceMap[serviceConfig.Name]; ok {
			err = oldService.Reload(i.ctx, serviceConfig)
			if err != nil {
				i.logger.Error().Str("service", serviceConfig.Name).Err(err).Msg("Error when updating services")
				return
			}
			newServiceMap[serviceConfig.Name] = oldService
		} else {
			newService := service.NewService(i.logger, serviceConfig)
			newService.UpdateRouter(i.router)
			err = newService.Start(i.ctx)
			if err != nil {
				i.logger.Error().Str("service", serviceConfig.Name).Err(err).Msg("Error when initializing services")
				return
			}
			newServiceMap[serviceConfig.Name] = newService
		}
	}

	i.outboundMap = newOutboundMap
	i.serviceMap = newServiceMap
}
