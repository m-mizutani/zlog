package zlog

import "reflect"

func (x *Logger) InjectInfra(infra *Infra) {
	x.infra = infra
}

func NewPurifier(filters Filters) *purifier {
	return newPurifier(filters)
}

func (x *purifier) Clone(v interface{}) interface{} {
	return x.clone(reflect.ValueOf(v), "").Interface()
}
