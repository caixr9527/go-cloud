package binding

import (
	"encoding/json"
	"net/http"
	"reflect"
)

type uriBinding struct {
	Params map[string]any
}

func (b uriBinding) Bind(r *http.Request, obj any) error {
	typeOf := reflect.TypeOf(obj)
	res := make(map[string]any)
	for i := 0; i < typeOf.NumField(); i++ {
		field := typeOf.Field(i)
		name := field.Name
		res[name] = b.Params[name]
	}
	marshal, _ := json.Marshal(res)
	json.Unmarshal(marshal, obj)
	return validate(obj)
}
