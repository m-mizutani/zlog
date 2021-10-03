package zlog

import "reflect"

type censor struct {
	filters Filters
}

func newCensor(filters Filters) *censor {
	return &censor{
		filters: filters,
	}
}

func (x *censor) clone(value reflect.Value, tag string) reflect.Value {
	adjustValue := func(ret reflect.Value) reflect.Value {
		switch value.Kind() {
		case reflect.Ptr, reflect.Map, reflect.Array, reflect.Slice:
			return ret
		default:
			return ret.Elem()
		}
	}

	src := value
	if value.Kind() == reflect.Ptr {
		src = src.Elem()
	}

	var dst reflect.Value
	if x.filters.IsSensitive(src.Interface(), tag) {
		dst = reflect.New(src.Type())
		if src.Kind() == reflect.String {
			dst.Elem().SetString(FilteredLabel)
		}
		return adjustValue(dst)
	}

	switch src.Kind() {
	case reflect.String:
		dst = reflect.New(src.Type())
		filtered := x.filters.ReplaceString(value.String())
		dst.Elem().SetString(filtered)

	case reflect.Struct:
		dst = reflect.New(src.Type())
		t := src.Type()

		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			fv := src.Field(i)
			if !fv.CanInterface() {
				continue
			}

			dst.Elem().Field(i).Set(x.clone(fv, f.Tag.Get("zlog")))
		}

	case reflect.Map:
		dst = reflect.MakeMap(src.Type())
		keys := src.MapKeys()
		for i := 0; i < src.Len(); i++ {
			mValue := src.MapIndex(keys[i])
			dst.SetMapIndex(keys[i], x.clone(mValue, ""))
		}

	case reflect.Array, reflect.Slice:
		dst = reflect.MakeSlice(src.Type(), src.Len(), src.Cap())
		for i := 0; i < src.Len(); i++ {
			dst.Index(i).Set(x.clone(src.Index(i), ""))
		}

	default:
		dst = reflect.New(src.Type())
		dst.Elem().Set(src)
	}

	return adjustValue(dst)
}
