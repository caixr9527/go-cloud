package factory

import (
	"github.com/caixr9527/go-cloud/common/utils"
	"github.com/caixr9527/go-cloud/component"
	"reflect"
)

func Get[T any](t T) T {
	name := utils.ObjName(t)
	val, ok := component.SinglePool.Get(name)
	if !ok {
		return t
	}
	return val.(T)
}

func GetSpecificField[T any](component any, field T) T {
	com := Get(component)
	if com == nil {
		return field
	}
	valueOf := reflect.ValueOf(com)
	numField := valueOf.NumField()
	for idx := 0; idx < numField; idx++ {
		val := valueOf.Field(idx)
		name := utils.ObjName(val)
		specificFieldName := utils.ObjName(field)
		if name == "" || specificFieldName == "" {
			continue
		} else if name == specificFieldName {
			return val.Interface().(T)
		}
	}
	return field
}

func GetField[T any](name string, t T) T {
	// todo 有问题需要修改
	get, ok := component.SinglePool.Get(name)
	if !ok {
		return t
	}
	valueOf := reflect.ValueOf(get)
	if valueOf.Kind() == reflect.Pointer {
		valueOf = reflect.ValueOf(valueOf.Elem())
	}
	numField := valueOf.NumField()
	for idx := 0; idx < numField; idx++ {
		field := valueOf.Field(idx)
		fieldName := field.String()
		//fieldName := utils.ObjName(field.Elem())
		specificFieldName := utils.ObjName(t)
		if fieldName == "" || specificFieldName == "" {
			continue
		} else if fieldName == specificFieldName {
			return field.Interface().(T)
		}
	}
	return t
}

func Create(obj any) {
	// todo zap.logger 有问题
	if obj == nil {
		return
	}
	v := reflect.ValueOf(obj)
	m := v.MethodByName("Name")
	var name string
	if m.IsValid() {
		call := m.Call(nil)
		if call != nil {
			name = call[0].String()
		}
	} else {
		name = utils.ObjName(obj)
	}

	component.SinglePool.Register(name, obj)
}

func Del(obj any) {
	name := utils.ObjName(obj)
	component.SinglePool.Del(name)
}
