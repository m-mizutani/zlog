package filter

import "strings"

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

type FieldPrefixFilter struct {
	prefix string
}

func FieldPrefix(prefix string) *FieldPrefixFilter {
	return &FieldPrefixFilter{
		prefix: prefix,
	}
}

func (x *FieldPrefixFilter) ReplaceString(s string) string {
	return s
}

func (x *FieldPrefixFilter) ShouldMask(fieldName string, value interface{}, tag string) bool {
	return strings.HasPrefix(fieldName, x.prefix)
}
