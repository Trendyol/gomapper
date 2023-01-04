package gomapper

import "reflect"

func isAnyNil(x any) bool {
	if x == nil {
		return true
	}

	switch reflect.TypeOf(x).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(x).IsNil()
	}

	return false
}

func isReflectValNil(value reflect.Value) bool {
	return value.Kind() == reflect.Ptr && value.IsNil()
}
