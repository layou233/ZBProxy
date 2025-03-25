package adapter

import (
	"fmt"

	"github.com/layou233/zbproxy/v3/common/set"
	"github.com/layou233/zbproxy/v3/config"
)

type OutboundReloadOptions struct {
	Router Router
	Config *config.Outbound
	Lists  map[string]set.StringSet
}

var _ RouteResourceProvider = (*OutboundReloadOptions)(nil)

func (o *OutboundReloadOptions) FindOutboundByName(name string) (Outbound, error) {
	// TODO: also implement independent outbound finder
	return o.Router.FindOutboundByName(name)
}

func (o *OutboundReloadOptions) FindListsByTag(tags []string) ([]set.StringSet, error) {
	if o.Lists != nil {
		lists := make([]set.StringSet, 0, len(tags))
		for _, tag := range tags {
			list, ok := o.Lists[tag]
			if !ok {
				return nil, fmt.Errorf("list not found [%s]", tag)
			}
			lists = append(lists, list)
		}
		return lists, nil
	}
	return o.Router.FindListsByTag(tags)
}
