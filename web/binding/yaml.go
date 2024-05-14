package binding

import (
	"errors"
	"github.com/caixr9527/go-cloud/web/validator"
	"gopkg.in/yaml.v3"
	"net/http"
)

type yamlBinding struct {
}

func (b yamlBinding) Bind(r *http.Request, obj any) error {
	body := r.Body
	if body == nil {
		return errors.New("invalid request")
	}
	decoder := yaml.NewDecoder(body)
	err := decoder.Decode(obj)
	if err != nil {
		return err
	}
	return validator.New().Struct(obj)
}
