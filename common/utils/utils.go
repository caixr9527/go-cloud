package utils

import "reflect"

func ObjName(t any) string {
	typeOf := reflect.TypeOf(t)
	kind := typeOf.Kind()
	name := typeOf.String()
	if kind == reflect.Pointer {
		name = typeOf.Elem().String()
	}
	return name
}
