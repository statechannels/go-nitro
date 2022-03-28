package helpers // import "github.com/statechannels/go-nitro/helpers"

import (
	"fmt"
	"reflect"
)

// HasShallowCopy takes two structs and checks if one is a shallow copy of another.
// This is done by using reflection and checking that reference types like slices and maps don't point to the same values.
func HasShallowCopy(a interface{}, b interface{}) bool {

	// Define a recursive function that can be used to determine if there exists a shallow copy between a and b
	var isShallow func(a, b reflect.Value) bool
	isShallow = func(aReflect reflect.Value, bReflect reflect.Value) bool {
		if aReflect.Kind() != bReflect.Kind() {
			panic(fmt.Errorf("cannot compare two different kinds %s and %s", aReflect.Kind(), bReflect.Kind()))
		}
		switch kind := aReflect.Kind(); kind {

		case reflect.Struct:
			for i := 0; i < aReflect.NumField(); i++ {
				aVal := aReflect.Field(i)
				bVal := bReflect.Field(i)
				if isShallow(aVal, bVal) {
					return true
				}
			}

		case reflect.Slice:
			for j := 0; j < aReflect.Len(); j++ {
				aVal := aReflect.Index(j)
				bVal := bReflect.Index(j)
				if isShallow(aVal, bVal) {
					return true
				}
			}
		case reflect.Map:
			itr := aReflect.MapRange()
			for itr.Next() {
				key := itr.Key()
				aVal := aReflect.MapIndex(key)
				bVal := bReflect.MapIndex(key)
				if isShallow(aVal, bVal) {
					return true
				}

			}
		case reflect.Ptr:
			return aReflect.Pointer() == bReflect.Pointer()
		default:
			return false

		}
		return false
	}
	aReflect := reflect.ValueOf(a)
	bReflect := reflect.ValueOf(b)
	// IF we have a pointer we're dealing with an interface so we unwrap it.
	if aReflect.Kind() == reflect.Ptr {

		aReflect = reflect.ValueOf(a).Elem()
		bReflect = reflect.ValueOf(b).Elem()
	}
	return isShallow(aReflect, bReflect)
}
