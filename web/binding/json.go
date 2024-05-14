package binding

import (
	"encoding/json"
	"errors"
	"github.com/caixr9527/go-cloud/web/validator"
	"net/http"
)

type jsonBinding struct {
	DisallowUnknownFields bool
	IsValidate            bool
}

func (b jsonBinding) Bind(r *http.Request, obj any) error {
	body := r.Body
	if body == nil {
		return errors.New("invalid request")
	}
	decoder := json.NewDecoder(body)
	if b.DisallowUnknownFields {
		decoder.DisallowUnknownFields()
	}
	err := decoder.Decode(obj)
	if err != nil {
		return err
	}
	return validator.New().Struct(obj)
}
