package binding

import (
	"encoding/xml"
	"errors"
	"github.com/caixr9527/go-cloud/web/validator"
	"net/http"
)

type xmlBinding struct {
}

func (b xmlBinding) Bind(r *http.Request, obj any) error {
	body := r.Body
	if body == nil {
		return errors.New("invalid request")
	}
	decoder := xml.NewDecoder(body)
	err := decoder.Decode(obj)
	if err != nil {
		return err
	}
	return validator.New().Struct(obj)
}
