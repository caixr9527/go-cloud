package factory

import (
	"errors"
	"github.com/caixr9527/go-cloud/component"
	"reflect"
)

func Get[T any](t T) (T, error) {
	typeOf := reflect.TypeOf(t)
	kind := typeOf.Kind()
	if kind == reflect.Pointer {
		return t, errors.New("cannot be a pointer type")
	}
	return component.SinglePool.Get(typeOf.Name()).(T), nil
}
