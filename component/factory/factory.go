package factory

import (
	"github.com/caixr9527/go-cloud/component"
	"reflect"
)

func Get[T any](t T) T {
	typeOf := reflect.TypeOf(t)
	kind := typeOf.Kind()
	name := typeOf.String()
	if kind == reflect.Pointer {
		name = typeOf.Elem().String()
	}
	return component.SinglePool.Get(name).(T)
}
