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

type Option struct {
	// If this is false(default); It doesn't fail
	// when the destination type contains fields not supplied by the source.
	// If this is true; All fields in the
	// destination object must exist in the source object.
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
			// If source is nil, it maps to an empty struct.
			sourceVal = reflect.New(sourceVal.Type().Elem())
		} else if sourceVal.Kind() == reflect.Ptr {
			sourceVal = sourceVal.Elem()
		}

		for i := 0; i < destVal.NumField(); i++ {
			if err := mapField(sourceVal, destVal, i, loose); err != nil {
				if !loose {
					return err
				}
			}
		}
	} else if destVal.Kind() == reflect.Slice {
		if isReflectValNil(sourceVal) {
			return nil
		} else if sourceVal.Kind() == reflect.Ptr {
			sourceVal = sourceVal.Elem()
		}
		return mapSlice(sourceVal, destVal, loose)
	} else if destVal.Kind() == reflect.Map {
		if isReflectValNil(sourceVal) {
			return nil
		} else if sourceVal.Kind() == reflect.Ptr {
			sourceVal = sourceVal.Elem()
		}
		return mapMap(sourceVal, destVal, loose)
	} else {
		return errors.New("error mapping values: currently not supported")
	}

	return nil
}

func mapSlice(sourceVal, destVal reflect.Value, loose bool) error {
	destType := destVal.Type()
	length := sourceVal.Len()
	target := reflect.MakeSlice(destType, length, length)

	for i := 0; i < length; i++ {
		val := reflect.New(destType.Elem()).Elem()
		if err := mapValues(sourceVal.Index(i), val, loose); err != nil {
			return err
		}
		target.Index(i).Set(val)
	}

	if length == 0 {
		if err := verifySliceTypesAreCompatible(sourceVal, destVal, loose); err != nil {
			return err
		}
	}

	destVal.Set(target)
	return nil
}

func mapMap(sourceVal, destVal reflect.Value, loose bool) error {
	// // Kaynak ve hedef verinin türlerini alın
	// sourceType := sourceVal.Type()
	// destType := destVal.Type()

	// // Kaynak ve hedef verinin eleman türlerini alın
	// sourceElemType := sourceType.Elem()
	// destElemType := destType.Elem()

	// // Eğer kaynak ve hedef verinin eleman türleri aynı ise, hedef veriyi kaynak veriden kopyalayın
	// if sourceElemType == destElemType {
	// 	destVal.Set(sourceVal)
	// 	return nil
	// }

	// // Eğer kaynak ve hedef verinin eleman türleri farklı ise, hedef veriyi oluşturun
	// destMap := reflect.MakeMap(destType)

	// // Kaynak verinin elemanlarını döngüyle gezin
	// for _, key := range sourceVal.MapKeys() {
	// 	// Kaynak verinin elemanını alın
	// 	sourceElem := sourceVal.MapIndex(key)

	// 	// Eğer hedef verinin eleman türü bir slice (dizi) ise, kaynak verinin elemanını hedef verinin elemanına dönüştürün
	// 	if destElemType.Kind() == reflect.Slice {
	// 		destElem := reflect.New(destElemType).Elem()
	// 		if err := mapSlice(sourceElem, destElem, loose); err != nil {
	// 			return err
	// 		}
	// 		destMap.SetMapIndex(key, destElem)
	// 	} else {
	// 		// Eğer hedef verinin eleman türü bir slice değilse, hedef verinin elemanını oluşturun ve kaynak verinin elemanını hedef verinin elemanına dönüştürün
	// 		destElem := reflect.New(destElemType).Elem()
	// 		if err := mapValues(sourceElem, destElem, loose); err != nil {
	// 			return err
	// 		}
	// 		destMap.SetMapIndex(key, destElem)
	// 	}
	// }

	// // Oluşturulan hedef veriyi, hedef veri değişkenine atayın
	// destVal.Set(destMap)

	return nil
}

// func verifyMapTypesAreCompatible(sourceVal, destVal reflect.Value, loose bool) error {
// 	dummyDest := reflect.New(reflect.PtrTo(destVal.Type()))
// 	dummySource := reflect.MakeMapWithSize(sourceVal.Type(), 1)
// 	return mapValues(dummySource, dummyDest.Elem(), loose)
// }

func verifySliceTypesAreCompatible(sourceVal, destVal reflect.Value, loose bool) error {
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
		if isReflectValNil(parentField) {
			return true
		}
	}
	return false
}
