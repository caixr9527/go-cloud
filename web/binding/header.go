package binding

import (
	"encoding/json"
	"net/http"
	"reflect"
)

type headerBinding struct {
}

func (b headerBinding) Bind(r *http.Request, obj any) error {
	typeOf := reflect.TypeOf(obj)
	res := make(map[string]any)
	for i := 0; i < typeOf.NumField(); i++ {
		field := typeOf.Field(i)
		name := field.Name
		res[name] = r.Header.Get(name)
	}
	marshal, _ := json.Marshal(res)
	json.Unmarshal(marshal, obj)
	return validate(obj)
}