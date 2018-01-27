package adapters

import (
	"regexp"
	"fmt"
)

type PortDescription interface {
	GetPrefix() string
	IsMulti() bool
	CheckTag(tag string) bool
}

func newSinglePortDescription(name string) *portDescription {
	return &portDescription{
		prefix:name,
		isMulti:false,
		matcher: func(tag string) bool {
			return tag == name
		},
	}
}

func newMultiPortDescription(name string) *portDescription {
	re, _ := regexp.Compile(fmt.Sprintf("^%s_[0-9]+$", name))
	return &portDescription{
		prefix:name,
		isMulti:true,
		matcher: func(tag string) bool {
			return re.Match([]byte(name))
		},
	}
}

type portDescription struct {
	prefix string
	isMulti bool
	matcher func(string) bool
}

func (d *portDescription) GetPrefix() string {
	return d.prefix
}

func (d *portDescription) IsMulti() bool {
	return d.isMulti
}

func (d *portDescription) CheckTag(tag string) bool {
	return d.matcher(tag)
}

