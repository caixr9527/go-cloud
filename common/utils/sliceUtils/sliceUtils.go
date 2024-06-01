package sliceUtils

import "reflect"

func IsSliceOrArray(data any) bool {
	v := reflect.ValueOf(data)
	return v.Kind() == reflect.Array || v.Kind() == reflect.Slice
}

func IsEmpty(data any) bool {
	if data == nil {
		return true
	}
	v := reflect.ValueOf(data)
	if !IsSliceOrArray(data) {
		return true
	} else {
		slice := v.Interface().([]interface{})
		return len(slice) == 0
	}
}

func IsNotEmpty(data any) bool {
	return !IsEmpty(data)
}

func Contains(slice any, value any) bool {
	if IsEmpty(slice) {
		return false
	}
	v := reflect.ValueOf(slice)
	data := v.Interface().([]interface{})
	for _, val := range data {
		if val == value {
			return true
		}
	}
	return false
}
