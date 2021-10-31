package filter

type TagFilter struct {
	SecureTags []string
}

const defaultFilterTagName = "secret"

func Tag(tags ...string) *TagFilter {
	if len(tags) == 0 {
		tags = []string{defaultFilterTagName}
	}
	return &TagFilter{
		SecureTags: tags,
	}
}

func (x *TagFilter) ReplaceString(s string) string { return s }

func (x *TagFilter) ShouldMask(fieldName string, value interface{}, tag string) bool {
	for i := range x.SecureTags {
		if x.SecureTags[i] == tag {
			return true
		}
	}
	return false
}
