// Package gomapper provides support for mapping between two different types
// with compatible fields. The intended application for this is when you use
// one set of types to represent DTOs (data transfer objects, e.g. json data),
// and a different set of types internally in the application. Using this
// package can help converting from one type to another.
//
// This package uses reflection to perform mapping which should be fine for
// all but the most demanding applications.
package gomapper

import (
	"fmt"
	"reflect"

	"errors"
)

type MapOptions struct {
	// If this is false(default); It doesn't fail
	// when the destination type contains fields not supplied by the source.
	// If this is true; All fields in the
	// destination object must exist in the source object.
	// Object hierarchies with nested structs and slices are supported, as long as
	// type types of nested structs/slices follow the same rules, i.e. all fields
	// in destination structs must be found on the source struct.
	Exact bool
}

func getDefaultMapOptions() *MapOptions {
	return &MapOptions{
		Exact: false,
	}
}

// Map uses parametric options to fill out the fields in dest with values from source.
// If options does not provided it uses default map options.
// Embedded/anonymous structs are supported.
// Values that are not exported/not public will not be mapped.
func Map(source, dest any, opts ...*MapOptions) error {
	mapOptions, err := validateMapOptions(opts...)
	if err != nil {
		return err
	}

	if isAnyNil(source) {
		return errors.New("source must not be nil")
	}

	if isAnyNil(dest) {
		return errors.New("dest must not be nil")
	}

	if reflect.TypeOf(dest).Kind() != reflect.Ptr {
		return errors.New("dest must be a pointer type")
	}

	sourceVal := reflect.ValueOf(source)

	if reflect.TypeOf(source).Kind() == reflect.Ptr {
		sourceVal = reflect.ValueOf(source).Elem()
	}

	var destVal = reflect.ValueOf(dest).Elem()

	return mapValues(sourceVal, destVal, !mapOptions.Exact)
}

func validateMapOptions(opts ...*MapOptions) (*MapOptions, error) {
	if len(opts) > 1 {
		return nil, errors.New("function accepts only one option as a parameter")
	}

	var mapOptions *MapOptions

	if len(opts) == 0 {
		mapOptions = getDefaultMapOptions()
	} else {
		mapOptions = opts[0]
	}

	return mapOptions, nil
}

func mapValues(sourceVal, destVal reflect.Value, loose bool) error {
	destType := destVal.Type()

	// If the types are equal, map to destination from the top.
	// This can cause side effects, because pointer fields will point
	// to the same structure. In practice we are using this tool for transfering
	// data between layers. Not using for deep copy purposes. This is acceptable.
	if destVal.CanSet() && destType == sourceVal.Type() {
		destVal.Set(sourceVal)

		return nil
	} else if destType.Kind() == reflect.Struct {
		if sourceVal.Type().Kind() == reflect.Ptr {
			if sourceVal.IsNil() {
				// If source is nil, it maps to an empty struct.
				sourceVal = reflect.New(sourceVal.Type().Elem())
			}
			sourceVal = sourceVal.Elem()
		}
		for i := 0; i < destVal.NumField(); i++ {
			if err := mapField(sourceVal, destVal, i, loose); err != nil {
				if !loose {
					return err
				}
			}
		}

		return nil
	} else if destType.Kind() == reflect.Ptr {
		if reflectValueIsNil(sourceVal) {
			return nil
		}
		val := reflect.New(destType.Elem())
		if err := mapValues(sourceVal, val.Elem(), loose); err != nil {
			return err
		}
		destVal.Set(val)

		return nil
	} else if destType.Kind() == reflect.Slice {
		return mapSlice(sourceVal, destVal, loose)
	} else {
		return errors.New("error mapping values: currently not supported")
	}
}

func mapSlice(sourceVal, destVal reflect.Value, loose bool) error {
	destType := destVal.Type()
	length := sourceVal.Len()
	target := reflect.MakeSlice(destType, length, length)
	for j := 0; j < length; j++ {
		val := reflect.New(destType.Elem()).Elem()
		if err := mapValues(sourceVal.Index(j), val, loose); err != nil {
			return err
		}
		target.Index(j).Set(val)
	}

	if length == 0 {
		if err := verifyArrayTypesAreCompatible(sourceVal, destVal, loose); err != nil {
			return err
		}
	}
	destVal.Set(target)

	return nil
}

func verifyArrayTypesAreCompatible(sourceVal, destVal reflect.Value, loose bool) error {
	dummyDest := reflect.New(reflect.PtrTo(destVal.Type()))
	dummySource := reflect.MakeSlice(sourceVal.Type(), 1, 1)
	return mapValues(dummySource, dummyDest.Elem(), loose)
}

func mapField(source, destVal reflect.Value, i int, loose bool) error {
	destType := destVal.Type()
	fieldName := destType.Field(i).Name

	defer func() {
		if r := recover(); r != nil {
			panic(fmt.Sprintf("error mapping field: %s. DestType: %v SourceType: %v Error: %v",
				fieldName, destType, source.Type(), r))
		}
	}()

	destField := destVal.Field(i)

	if !destField.CanSet() {
		if loose {
			return nil
		} else {
			return fmt.Errorf("error mapping field: %s. Field can not set! DestType: %v SourceType: %v",
				fieldName, destType, source.Type())
		}
	}

	if destType.Field(i).Anonymous {
		return mapValues(source, destField, loose)
	} else {
		if valueIsContainedInNilEmbeddedType(source, fieldName) {
			return nil
		}

		sourceField := source.FieldByName(fieldName)
		if (sourceField == reflect.Value{}) {
			if loose {
				return nil
			}

			return fmt.Errorf("error mapping field: %s. SourceType: %v does not contain related field. DestType: %v",
				fieldName, source.Type(), destType)
		}

		return mapValues(sourceField, destField, loose)
	}
}

func valueIsContainedInNilEmbeddedType(source reflect.Value, fieldName string) bool {
	structField, _ := source.Type().FieldByName(fieldName)
	ix := structField.Index
	if len(structField.Index) > 1 {
		parentField := source.FieldByIndex(ix[:len(ix)-1])
		if reflectValueIsNil(parentField) {
			return true
		}
	}
	return false
}
