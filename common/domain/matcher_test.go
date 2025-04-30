package domain

import "testing"

func TestMatcher(t *testing.T) {
	matcher := NewMatcher(
		[]string{ // domain
			"eXample.com", "example.com.", "example.org",
		}, []string{ // domain suffix
			"example.net", ".exampLe.invalid",
		})
	if !matcher.Match("exaMple.com") {
		t.Error("example.com is not matched")
	}
	if !matcher.Match("example.com.") {
		t.Error("example.com. is not matched")
	}
	if !matcher.Match("example.org") {
		t.Error("example.org is not matched")
	}
	if matcher.Match("sub.example.org") {
		t.Error("sub.example.org is matched")
	}
	if !matcher.Match("example.net") {
		t.Error("example.net is not matched")
	}
	if !matcher.Match("any.example.net") {
		t.Error("any.example.net is not matched")
	}
	if !matcher.Match("any.one.example.net") {
		t.Error("any.one.one.example.net is not matched")
	}
	if matcher.Match("example.invAlid") {
		t.Error("example.invalid is matched")
	}
	if !matcher.Match("any.example.invalid") {
		t.Error("any.example.invalid is not matched")
	}
}
