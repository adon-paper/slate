package helper

import (
	"reflect"
	"strings"
)

func Empty(a interface{}) bool {
	if a == nil {
		return true
	}

	kind := reflect.TypeOf(a).Kind()
	value := reflect.ValueOf(a)

	switch kind {
	case reflect.String:
		return value.String() == ""
	case reflect.Int, reflect.Int32, reflect.Int64, reflect.Int8, reflect.Int16, reflect.Uint, reflect.Uint32, reflect.Uint64, reflect.Uint8, reflect.Uint16:
		return value.Int() == 0
	case reflect.Float32, reflect.Float64:
		return value.Float() == 0
	case reflect.Struct:
		return value.IsZero()
	case reflect.Ptr:
		return value.IsNil()
	case reflect.Slice, reflect.Array, reflect.Map:
		return value.Len() == 0
	case reflect.Bool:
		return !value.Bool()
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

func IsAggregates(s string) bool {
	aggregate := []string{"COUNT", "SUM"}
	if stringContainInSlice(s, aggregate) {
		return true
	}
	return false
}

func stringContainInSlice(a string, list []string) bool {
	for _, b := range list {
		if strings.Contains(strings.ToLower(a), strings.ToLower(b)) {
			return true
		}
	}
	return false
}
