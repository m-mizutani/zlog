package filter

type TagFilter struct {
	SecureTag string
}

func Tag() *TagFilter {
	return &TagFilter{
		SecureTag: "secure",
	}
}

func (x *TagFilter) ReplaceString(s string) string { return s }

func (x *TagFilter) IsSensitive(value interface{}, tag string) bool {
	return x.SecureTag == tag
}
