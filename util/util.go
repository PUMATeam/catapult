package util

import (
	"reflect"
)

// StructToMap converts a struct to a map
func StructToMap(in interface{}, f func(s string) string) map[string]string {
	out := make(map[string]string)
	v := reflect.ValueOf(in)
	typ := v.Type()

	for i := 0; i < v.NumField(); i++ {
		// gets us a StructField
		fi := f(typ.Field(i).Name)
		// set key of map to value in struct field
		out[fi] = v.Field(i).String()
	}

	return out
}
