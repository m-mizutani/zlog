package zlog

func (x *Logger) InjectInfra(infra *Infra) {
	x.infra = infra
}
