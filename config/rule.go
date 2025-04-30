package config

import (
	"encoding/json"

	"github.com/layou233/zbproxy/v3/common/jsonx"
)

type _RuleParameterListableString struct {
	Lists     jsonx.Listable[string]
	LowerCase bool
}

type RuleParameterListableString _RuleParameterListableString

var _ json.Unmarshaler = (*RuleParameterListableString)(nil)

func (r *RuleParameterListableString) UnmarshalJSON(bytes []byte) error {
	err := json.Unmarshal(bytes, (*_RuleParameterListableString)(r))
	if err == nil {
		return nil
	}
	err = json.Unmarshal(bytes, &r.Lists)
	if err == nil {
		r.LowerCase = false
		return nil
	}
	return err
}
