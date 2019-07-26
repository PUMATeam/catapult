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
		tag := f(typ.Field(i).Tag.Get("json"))

		if tag == "ignore" {
			continue
		}

		fi := f(typ.Field(i).Name)
		// set key of map to value in struct field
		key := fi
		if tag != "" {
			key = tag
		}
		out[key] = v.Field(i).String()
	}

	return out
}
