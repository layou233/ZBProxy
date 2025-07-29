package config

import (
	"encoding/json"

	"github.com/layou233/zbproxy/v3/common/jsonx"
)

type Router struct {
	DefaultOutbound string  `json:",omitempty"`
	Rules           []*Rule `json:",omitempty"`
}

type Rule struct {
	Type      string
	Parameter json.RawMessage `json:",omitempty"`
	//SubRules []Rule `json:",omitempty"`
	Rewrite  RuleRewrite            `json:",omitempty"`
	Sniff    jsonx.Listable[string] `json:",omitempty"`
	Outbound string                 `json:",omitempty"`
	Invert   bool                   `json:",omitempty"`
}

type RuleRewrite struct {
	TargetAddress string                `json:",omitempty"`
	TargetPort    uint16                `json:",omitempty"`
	Minecraft     *ruleRewriteMinecraft `json:",omitempty"`
}

type ruleRewriteMinecraft struct {
	Hostname string `json:",omitempty"`
	Port     uint16 `json:",omitempty"`
	Intent   int8   `json:",omitempty"`
}

type RuleDomain struct {
	Domain       jsonx.Listable[string] `json:",omitempty"`
	DomainSuffix jsonx.Listable[string] `json:",omitempty"`
}
