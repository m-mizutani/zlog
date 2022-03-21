package zlog

import "reflect"

func NewMasking(filters Filters) *masking {
	return newMasking(filters)
}

func (x *masking) Clone(v interface{}) interface{} {
	return x.clone("", reflect.ValueOf(v), "").Interface()
}
