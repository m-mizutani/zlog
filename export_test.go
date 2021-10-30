package zlog

import "reflect"

func (x *Logger) InjectInfra(infra *Infra) {
	x.infra = infra
}

func NewMasking(filters Filters) *masking {
	return newMasking(filters)
}

func (x *masking) Clone(v interface{}) interface{} {
	return x.clone("", reflect.ValueOf(v), "").Interface()
}
