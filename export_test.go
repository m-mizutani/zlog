package zlog

import "reflect"

func (x *Logger) InjectInfra(infra *Infra) {
	x.infra = infra
}

func NewCensor(filters Filters) *censor {
	return newCensor(filters)
}

func (x *censor) Clone(v interface{}) interface{} {
	return x.clone(reflect.ValueOf(v), "").Interface()
}
