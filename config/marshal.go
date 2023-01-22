package config

import (
	"encoding/json"
	"github.com/layou233/ZBProxy/common/set"
)

type configMainTemp struct {
	Services []*ConfigProxyService
	Lists    map[string][]string
}

func (c *configMain) MarshalJSON() ([]byte, error) {
	var list map[string][]string
	if l := len(c.Lists); l == 0 { // if nothing in Lists
		return json.Marshal( // empty map
			configMainTemp{
				Services: c.Services,
				Lists:    nil,
			},
		)
	} else {
		list = make(map[string][]string, l) // map size init
		for k, v := range c.Lists {
			for k1, _ := range *v {
				list[k] = append(list[k], k1)
			}
		}
	}
	return json.Marshal(
		configMainTemp{
			Services: c.Services,
			Lists:    list,
		},
	)
}

func (c *configMain) UnmarshalJSON(data []byte) (err error) {
	var configTemp configMainTemp
	err = json.Unmarshal(data, &configTemp)
	if err != nil {
		return err
	}
	// log.Println("Lists:", configTemp.Lists)
	if l := len(configTemp.Lists); l == 0 { // if nothing in Lists
		c.Lists = map[string]*set.StringSet{} // empty map
	} else {
		c.Lists = make(map[string]*set.StringSet, l) // map size init
		for k, v := range configTemp.Lists {
			// log.Println("List: Loading", k, "value:", v)
			list := set.NewStringSetFromSlice(v)
			c.Lists[k] = &list
		}
	}
	configTemp.Lists = nil // free memory
	c.Services = configTemp.Services
	return err
}
