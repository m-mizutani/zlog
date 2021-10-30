package zlog

type Filter interface {
	ReplaceString(s string) string
	ShouldMask(fieldName string, value interface{}, tag string) bool
}

type Filters []Filter

func (x Filters) ReplaceString(s string) string {
	for _, f := range x {
		s = f.ReplaceString(s)
	}
	return s
}

func (x Filters) ShouldMask(fieldName string, value interface{}, tag string) bool {
	for _, f := range x {
		if f.ShouldMask(fieldName, value, tag) {
			return true
		}
	}
	return false
}
