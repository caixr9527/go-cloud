package factory

import (
	"github.com/caixr9527/go-cloud/common"
	"github.com/caixr9527/go-cloud/common/utils"
	"github.com/caixr9527/go-cloud/internal/component"
	"reflect"
)

func Get[T any](t T) T {
	name := getName(t)
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
	name := getName(obj)

	component.SinglePool.Register(name, obj)
}

func getName(obj any) string {
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
	return name
}

func Del(obj any) {
	name := utils.ObjName(obj)
	component.SinglePool.Del(name)
}

func RegisterComponent(c ...common.Bean) {
	component.Components = append(component.Components, c...)
}

func RegisterBean(b ...common.Bean) {
	component.Beans = append(component.Beans, b...)
}
