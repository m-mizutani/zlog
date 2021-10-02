package zlog

import "time"

type Infra struct {
	Now func() time.Time
}

func newInfra() *Infra {
	return &Infra{
		Now: time.Now,
	}
}
