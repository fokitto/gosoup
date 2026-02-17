package gosoup

import (
	"regexp"
	"strings"
)

type Predicate func(*Tag) bool

func HasName(name string) Predicate {
	return func(tag *Tag) bool {
		return tag.Name == name
	}
}

func HasAttr(attr string) Predicate {
	return func(tag *Tag) bool {
		_, ok := tag.Attrs[attr]
		return ok
	}
}

func HasNoAttr(attr string) Predicate {
	return func(tag *Tag) bool {
		_, ok := tag.Attrs[attr]
		return !ok
	}
}

func HasClass(class string) Predicate {
	return func(tag *Tag) bool {
		tagClass, ok := tag.Attrs["class"]
		if !ok {
			return false
		}
		for _, entry := range strings.Split(tagClass, " ") {
			if entry == class {
				return true
			}
		}
		return false
	}
}

func HasNoClass() Predicate {
	return func(tag *Tag) bool {
		return HasNoAttr("class")(tag)
	}
}

func AttrEq(attr string, value string) Predicate {
	return func(tag *Tag) bool {
		if tagAttr, ok := tag.Attrs[attr]; ok {
			return tagAttr == value
		}
		return false
	}
}

func AttrContains(attr string, substr string) Predicate {
	return func(tag *Tag) bool {
		if tagAttr, ok := tag.Attrs[attr]; ok {
			return strings.Contains(tagAttr, substr)
		}
		return false
	}
}

func AttrMatch(attr string, pattern *regexp.Regexp) Predicate {
	return func(tag *Tag) bool {
		if tagAttr, ok := tag.Attrs[attr]; ok {
			return pattern.MatchString(tagAttr)
		}
		return false
	}
}

func All(predicates ...Predicate) Predicate {
	return func(tag *Tag) bool {
		for _, predicate := range predicates {
			if !predicate(tag) {
				return false
			}
		}
		return true
	}
}

func Any(predicates ...Predicate) Predicate {
	return func(tag *Tag) bool {
		for _, predicate := range predicates {
			if predicate(tag) {
				return true
			}
		}
		return false
	}
}
