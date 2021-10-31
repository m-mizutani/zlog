package zlog

type Filter interface {
	ReplaceString(s string) string
	ShouldMask(fieldName string, value interface{}, tag string) bool
}

type Filters []Filter

// ReplaceString is called when checking string type. The argument is the value to be checked, and the return value should be the value to be replaced. If nothing needs to be done, the method should return the argument as is. This method is intended for the case where you want to hide a part of a string.
func (x Filters) ReplaceString(s string) string {
	for _, f := range x {
		s = f.ReplaceString(s)
	}
	return s
}

// ShouldMask is called for all values to be checked. The field name of the value to be checked, the value to be checked, and tag value if the structure has `zlog` tag will be passed as arguments. If the return value is false, nothing is done; if it is true, the entire field is hidden. All hidden values will be replaced with the value "[filtered]" of string type, regardless of their original type.
func (x Filters) ShouldMask(fieldName string, value interface{}, tag string) bool {
	for _, f := range x {
		if f.ShouldMask(fieldName, value, tag) {
			return true
		}
	}
	return false
}
