package zlog

type Filter interface {
	ReplaceString(s string) string
	ShouldMask(value interface{}, tag string) bool
}

type Filters []Filter

func (x Filters) ReplaceString(s string) string {
	for _, f := range x {
		s = f.ReplaceString(s)
	}
	return s
}

func (x Filters) ShouldMask(value interface{}, tag string) bool {
	for _, f := range x {
		if f.ShouldMask(value, tag) {
			return true
		}
	}
	return false
}
