package domain

import (
	"sort"
	"strings"
	"unicode/utf8"
)

type Matcher struct {
	set succinctSet
}

func NewMatcher(domains []string, domainSuffix []string) *Matcher {
	domainList := make([]string, 0, len(domains)+2*len(domainSuffix))
	seen := make(map[string]bool, len(domainList))
	for _, domain := range domainSuffix {
		domain = strings.ToLower(domain)
		if seen[domain] {
			continue
		}
		seen[domain] = true
		if domain[0] == '.' {
			domainList = append(domainList, reverseDomainSuffix(domain))
		} else {
			domainList = append(domainList, reverseDomainRoot(domain))
		}
	}
	for _, domain := range domains {
		domain = strings.ToLower(domain)
		if seen[domain] {
			continue
		}
		seen[domain] = true
		domainList = append(domainList, reverseDomain(domain))
	}
	sort.Strings(domainList)
	return &Matcher{newSuccinctSet(domainList)}
}

func (m *Matcher) Match(domain string) bool {
	return m.set.Has(reverseDomain(strings.ToLower(domain)))
}

func reverseDomain(domain string) string {
	l := len(domain)
	b := make([]byte, l)
	for i := 0; i < l; {
		r, n := utf8.DecodeRuneInString(domain[i:])
		i += n
		utf8.EncodeRune(b[l-i:], r)
	}
	return string(b)
}

func reverseDomainSuffix(domain string) string {
	l := len(domain)
	b := make([]byte, l+1)
	for i := 0; i < l; {
		r, n := utf8.DecodeRuneInString(domain[i:])
		i += n
		utf8.EncodeRune(b[l-i:], r)
	}
	b[l] = prefixLabel
	return string(b)
}

func reverseDomainRoot(domain string) string {
	l := len(domain)
	b := make([]byte, l+1)
	for i := 0; i < l; {
		r, n := utf8.DecodeRuneInString(domain[i:])
		i += n
		utf8.EncodeRune(b[l-i:], r)
	}
	b[l] = rootLabel
	return string(b)
}
