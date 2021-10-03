package filter

import (
	"reflect"

	"github.com/m-mizutani/zlog"
)

type TypeFilter struct {
	target reflect.Type
	zlog.Filter
}

func Type(t interface{}) *TypeFilter {
	return &TypeFilter{
		target: reflect.TypeOf(t),
	}
}

func (x *TypeFilter) ReplaceString(s string) string { return s }

func (x *TypeFilter) IsSensitive(value interface{}, tag string) bool {
	return x.target == reflect.TypeOf(value)
}
