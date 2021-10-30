package filter

import (
	"regexp"

	"github.com/m-mizutani/zlog"
)

type PhoneNumberFilter struct {
	RegexSet []regexp.Regexp
}

func PhoneNumber() *PhoneNumberFilter {
	return &PhoneNumberFilter{
		RegexSet: []regexp.Regexp{
			*regexp.MustCompile("[0-9]{2,4}-[0-9]{2,4}-[0-9]{4}"), // japan phone number format
		},
	}
}

func (x *PhoneNumberFilter) ReplaceString(s string) string {
	for _, p := range x.RegexSet {
		s = p.ReplaceAllString(s, zlog.FilteredLabel)
	}
	return s
}

func (x *PhoneNumberFilter) ShouldMask(fieldName string, value interface{}, tag string) bool {
	return false
}
