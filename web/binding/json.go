package binding

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/json-iterator/go/extra"
	"net/http"
)

type jsonBinding struct {
	DisallowUnknownFields bool
}

func (b jsonBinding) Bind(r *http.Request, obj any) error {
	body := r.Body
	if err := checkBody(body); err != nil {
		return err
	}
	extra.RegisterFuzzyDecoders()
	decoder := jsoniter.NewDecoder(body)
	if b.DisallowUnknownFields {
		decoder.DisallowUnknownFields()
	}
	if err := decoder.Decode(obj); err != nil {
		return err
	}
	return validate(obj)
}
