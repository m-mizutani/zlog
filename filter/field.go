package filter

type FieldFilter struct {
	target string
}

func Field(target string) *FieldFilter {
	return &FieldFilter{
		target: target,
	}
}

func (x *FieldFilter) ReplaceString(s string) string {
	return s
}

func (x *FieldFilter) ShouldMask(fieldName string, value interface{}, tag string) bool {
	return x.target == fieldName
}
