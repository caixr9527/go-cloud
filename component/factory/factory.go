package factory

import (
	"github.com/caixr9527/go-cloud/common/utils"
	"github.com/caixr9527/go-cloud/component"
	"reflect"
)

func Get[T any](t T) T {
	// todo fix
	name := utils.ObjName(t)
	val, ok := component.SinglePool.Get(name)
	if !ok {
		return t
	}
	return val.(T)
}

func Create(obj any) {
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
