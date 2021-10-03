package zlog

type Filter interface {
	ReplaceString(s string) string
	IsSensitive(value interface{}, tag string) bool
}

type Filters []Filter

func (x Filters) ReplaceString(s string) string {
	for _, f := range x {
		s = f.ReplaceString(s)
	}
	return s
}

func (x Filters) IsSensitive(value interface{}, tag string) bool {
	for _, f := range x {
		if f.IsSensitive(value, tag) {
			return true
		}
	}
	return false
}
