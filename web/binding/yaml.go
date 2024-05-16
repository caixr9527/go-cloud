package binding

import (
	"gopkg.in/yaml.v3"
	"net/http"
)

type yamlBinding struct {
}

func (b yamlBinding) Bind(r *http.Request, obj any) error {
	body := r.Body
	if err := checkBody(body); err != nil {
		return err
	}
	decoder := yaml.NewDecoder(body)
	err := decoder.Decode(obj)
	if err != nil {
		return err
	}
	return validate(obj)
}
