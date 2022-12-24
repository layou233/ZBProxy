package access

import (
	"fmt"

	"github.com/layou233/ZBProxy/common/set"
)

func GetTargetList(lists map[string]*set.StringSet, listName string) (*set.StringSet, error) {
	sets, ok := lists[listName]
	if ok {
		return sets, nil
	}
	return nil, fmt.Errorf("list %q not found", listName)
}
