package filter

import (
	"strings"

	"github.com/m-mizutani/zlog"
)

type ValueFilter struct {
	target string
}

func Value(target string) *ValueFilter {
	return &ValueFilter{
		target: target,
	}
}

func (x *ValueFilter) ReplaceString(s string) string {
	return strings.ReplaceAll(s, x.target, zlog.FilteredLabel)
}

func (x *ValueFilter) ShouldMask(value interface{}, tag string) bool {
	return false
}
