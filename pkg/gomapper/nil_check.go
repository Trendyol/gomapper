package gomapper

import "reflect"

func IsAnyNil(x any) bool {
	if x == nil {
		return true
	}

	switch reflect.TypeOf(x).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(x).IsNil()
	}

	return false
}

func ReflectValueIsNil(value reflect.Value) bool {
	return value.Type().Kind() == reflect.Ptr && value.IsNil()
}
