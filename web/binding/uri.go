package binding

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/json-iterator/go/extra"
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
	extra.RegisterFuzzyDecoders()
	var marshal []byte
	var err error
	if marshal, err = jsoniter.Marshal(res); err != nil {
		return err
	}
	if err = jsoniter.Unmarshal(marshal, obj); err != nil {
		return err
	}
	return validate(obj)
}
