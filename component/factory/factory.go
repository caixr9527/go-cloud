package factory

import (
	"github.com/caixr9527/go-cloud/component"
	"reflect"
)

func Get[T any](t T) T {
	typeOf := reflect.TypeOf(t)
	kind := typeOf.Kind()
	name := typeOf.Name()
	if kind == reflect.Pointer {
		name = typeOf.Elem().Name()
	}
	return component.SinglePool.Get(name).(T)
}
