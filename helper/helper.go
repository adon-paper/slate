package helper

import (
	"reflect"
	"time"
)

func Empty(a interface{}) bool {
	var emptyTime time.Time
	var emptyFloat64 float64
	var emptyFloat32 float32
	var emptyInt64 int64
	var emptyInt32 int32
	var emptyInt int

	if a == false ||
		a == "" ||
		a == emptyInt64 ||
		a == emptyInt32 ||
		a == emptyInt ||
		a == emptyFloat32 ||
		a == emptyFloat64 ||
		a == emptyTime ||
		a == nil {
		return true
	}

	switch reflect.TypeOf(a).Kind() {
	case reflect.Array, reflect.Slice:
		arr := reflect.ValueOf(a)
		if arr.Len() == 0 {
			return true
		}
	}
	return false
}

func MergeMaps(maps ...map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}

	return result
}
