package access

import (
	"fmt"
	"github.com/layou233/ZBProxy/common/set"
	"github.com/layou233/ZBProxy/config"
)

func IsStringInList(item string, stringSet *set.StringSet) bool {
	return stringSet.Has(item)
}

func GetTargetList(listName string) (*set.StringSet, error) {
	set, ok := config.Lists[listName]
	if ok {
		return set, nil
	}
	return nil, fmt.Errorf("list %q not found", listName)
}
