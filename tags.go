package hl7

import (
	"strconv"
	"strings"
)

type tagOptions string

func (o tagOptions) Required() bool {
	return o.Contains("required")
}

func (o tagOptions) Group() bool {
	return o.Contains("group")
}

func (o tagOptions) Optional() bool {
	return !o.Required()
}

func (o tagOptions) Contains(optionName string) bool {
	if len(o) == 0 {
		return false
	}
	s := string(o)
	for s != "" {
		var name string
		name, s, _ = strings.Cut(s, ",")
		if name == optionName {
			return true
		}
	}

	return false
}

type hl7Tag struct {
	Name    string
	Options tagOptions
}

func parseTag(tag string) hl7Tag {
	name, opt, found := strings.Cut(tag, ",")
	if !found {
		switch name {
		default:
			return hl7Tag{Name: name}
		case "group", "required":
			return hl7Tag{
				Options: tagOptions(name),
			}
		}
	}
	return hl7Tag{
		Name:    name,
		Options: tagOptions(opt),
	}
}

func (t hl7Tag) Index() (int, bool) {
	n, err := strconv.Atoi(t.Name)
	if err != nil {
		return 0, false
	}

	return n, true
}
