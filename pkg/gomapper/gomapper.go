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

type Option struct {
	// If this is false(default); It does not generate an error when the target type contains fields but these fields are not found in the source.
	// If this is true; All fields in the destination object must exist in the source object.
	// Also if this is true private destination fields must be supplied, that means if private destination field does not map automatically
	// from the upper object hierarchy then it will produce an error.
	// Object hierarchies with nested structs and slices are supported, as long as
	// type types of nested structs/slices follow the same rules, i.e. all fields
	// in destination structs must be found on the source struct.
	Exact bool
}

func getDefaultOption() *Option {
	return &Option{
		Exact: false,
	}
}

// Map uses parametric options to fill out the fields in dest with values from source.
// If options does not provided it uses default map options.
// Embedded/anonymous structs are supported.
// Values that are not exported/not public will not be mapped.
func Map(source, dest any, options ...*Option) error {
	option, err := verifyMapOption(options...)
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

	return mapValues(sourceVal, reflect.ValueOf(dest).Elem(), !option.Exact)
}

func verifyMapOption(options ...*Option) (*Option, error) {
	if len(options) > 1 {
		return nil, errors.New("only one option is accepted as a parameter")
	}

	var option *Option

	if len(options) == 0 {
		option = getDefaultOption()
	} else {
		option = options[0]
	}

	return option, nil
}

func mapValues(sourceVal, destVal reflect.Value, loose bool) error {
	// If the types are equal, map to destination from the top.
	// This can cause side effects, because pointer fields will point
	// to the same structure. In practice we are using this tool for transferring
	// data between layers. Not using for deep copy purposes. This is acceptable.
	if destVal.CanSet() && destVal.Type() == sourceVal.Type() {
		destVal.Set(sourceVal)
	} else if destVal.Kind() == reflect.Ptr {
		if isReflectValNil(sourceVal) {
			return nil
		}
		destValZeroPtr := reflect.New(destVal.Type().Elem())
		if err := mapValues(sourceVal, destValZeroPtr.Elem(), loose); err != nil {
			return err
		}
		destVal.Set(destValZeroPtr)
	} else if destVal.Kind() == reflect.Struct {
		if isReflectValNil(sourceVal) {
			// If source is nil, make a new default value of source's type.
			sourceVal = reflect.New(sourceVal.Type().Elem())
		}
		if sourceVal.Kind() == reflect.Ptr {
			sourceVal = sourceVal.Elem()
		}
		if sourceVal.Kind() != reflect.Struct {
			return errors.New("error mapping values: dest kind: struct, source kind: " + sourceVal.Kind().String())
		}
		for i := 0; i < destVal.NumField(); i++ {
			if err := mapField(sourceVal, destVal, i, loose); err != nil {
				return err
			}
		}
	} else if destVal.Kind() == reflect.Slice {
		if isReflectValNil(sourceVal) {
			return nil
		}
		if sourceVal.Kind() == reflect.Ptr {
			sourceVal = sourceVal.Elem()
		}
		if sourceVal.Kind() != reflect.Slice {
			return errors.New("error mapping values: dest kind: slice, source kind: " + sourceVal.Kind().String())
		}
		return mapSlice(sourceVal, destVal, loose)
	} else if destVal.Kind() == reflect.Map {
		if isReflectValNil(sourceVal) {
			return nil
		}
		if sourceVal.Kind() == reflect.Ptr {
			sourceVal = sourceVal.Elem()
		}
		if sourceVal.Kind() != reflect.Map {
			return errors.New("error mapping values: dest kind: map, source kind: " + sourceVal.Kind().String())
		}
		return mapMap(sourceVal, destVal, loose)
	} else {
		return errors.New(fmt.Sprintf("error mapping values: types are not compatible: Source Type: %s, Dest Type: %s", sourceVal.Type().Name(), destVal.Type().Name()))
	}

	return nil
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
			return errors.New(fmt.Sprintf("error mapping field: %s. Field can not set! DestType: %v SourceType: %v",
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

			return errors.New(fmt.Sprintf("error mapping field: %s. SourceType: %v does not contain related field. DestType: %v",
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
		if isReflectValNil(parentField) {
			return true
		}
	}
	return false
}

func mapSlice(sourceVal, destVal reflect.Value, loose bool) error {
	destType := destVal.Type()
	sourceLength := sourceVal.Len()
	target := reflect.MakeSlice(destType, sourceLength, sourceLength)

	for i := 0; i < sourceLength; i++ {
		val := reflect.New(destType.Elem()).Elem()
		if err := mapValues(sourceVal.Index(i), val, loose); err != nil {
			return err
		}
		target.Index(i).Set(val)
	}

	if sourceLength == 0 {
		if err := verifySliceTypesAreCompatible(sourceVal, destVal, loose); err != nil {
			return err
		}
	}

	destVal.Set(target)
	return nil
}

func verifySliceTypesAreCompatible(sourceVal, destVal reflect.Value, loose bool) error {
	dummyDest := reflect.New(reflect.PtrTo(destVal.Type())).Elem()
	dummySource := reflect.MakeSlice(sourceVal.Type(), 1, 1)
	return mapValues(dummySource, dummyDest, loose)
}

func mapMap(sourceVal, destVal reflect.Value, loose bool) error {
	sourceKeyType := sourceVal.Type().Key()
	destType := destVal.Type()
	destKeyType := destType.Key()

	if sourceKeyType.Name() != destKeyType.Name() {
		return errors.New(fmt.Sprintf("error mapping maps: map key types are not equal: Source Key Type: %s, Dest Key Type: %s", sourceKeyType.Name(), destKeyType.Name()))
	}

	sourceLength := sourceVal.Len()
	targetMap := reflect.MakeMapWithSize(destType, sourceLength)

	for _, key := range sourceVal.MapKeys() {
		sourceElem := sourceVal.MapIndex(key)

		destElem := reflect.New(destType.Elem()).Elem()
		if err := mapValues(sourceElem, destElem, loose); err != nil {
			return err
		}
		targetMap.SetMapIndex(key, destElem)
	}

	if sourceLength == 0 {
		if err := verifyMapElemTypesAreCompatible(sourceVal, destVal, loose); err != nil {
			return err
		}
	}

	destVal.Set(targetMap)
	return nil
}

func verifyMapElemTypesAreCompatible(sourceVal, destVal reflect.Value, loose bool) error {
	dummyDestElem := reflect.New(destVal.Type().Elem()).Elem()
	dummySourceElem := reflect.New(sourceVal.Type().Elem()).Elem()
	return mapValues(dummySourceElem, dummyDestElem, loose)
}
