package set

type StringSet map[string]struct{}

func (s StringSet) Has(item string) (ok bool) {
	_, ok = s[item]
	return
}

func (s StringSet) Add(item string) {
	s[item] = struct{}{}
}

func (s StringSet) Delete(item string) {
	delete(s, item)
}

func NewStringSetFromSlice(slice []string) StringSet {
	s := make(StringSet, len(slice))
	for _, item := range slice {
		s.Add(item)
	}
	return s
}
