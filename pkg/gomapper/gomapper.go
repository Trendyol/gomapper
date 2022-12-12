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

	"github.com/pkg/errors"
)

// Map fills out the fields in dest with values from source. All fields in the
// destination object must exist in the source object.
//
// Object hierarchies with nested structs and slices are supported, as long as
// type types of nested structs/slices follow the same rules, i.e. all fields
// in destination structs must be found on the source struct.
//
// Embedded/anonymous structs are supported
//
// Values that are not exported/not public will not be mapped.
//
// It is a design decision to panic when a field cannot be mapped in the
// destination to ensure that a renamed field in either the source or
// destination does not result in subtle silent bug.
func Map(source, dest any) error {
	return mapCommon(source, dest, false)
}

// MapLoose works just like Map, except it doesn't fail when the destination
// type contains fields not supplied by the source.
//
// This function is meant to be a temporary solution - the general idea is
// that the Map function should take a number of options that can modify its
// behavior - but I'd rather not add that functionality before I have a better
// idea what is a good options format.
func MapLoose(source, dest any) error {
	return mapCommon(source, dest, true)
}

func mapCommon(source, dest any, loose bool) error {
	if IsAnyNil(source) {
		return errors.New("source must not be nil")
	}

	if IsAnyNil(dest) {
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

	return mapValues(sourceVal, destVal, loose)
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
		if ReflectValueIsNil(sourceVal) {
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
			return errors.New(
				fmt.Sprintf("error mapping field: %s. Field can not set! DestType: %v SourceType: %v",
					fieldName, destType, source.Type()))
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

			return errors.New(
				fmt.Sprintf("error mapping field: %s. SourceType: %v does not contain related field. DestType: %v",
					fieldName, source.Type(), destType))
		}

		return mapValues(sourceField, destField, loose)
	}
}

func valueIsContainedInNilEmbeddedType(source reflect.Value, fieldName string) bool {
	structField, _ := source.Type().FieldByName(fieldName)
	ix := structField.Index
	if len(structField.Index) > 1 {
		parentField := source.FieldByIndex(ix[:len(ix)-1])
		if ReflectValueIsNil(parentField) {
			return true
		}
	}
	return false
}
