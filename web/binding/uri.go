package binding

import (
	"net/http"
)

type uriBinding struct {
	Params map[string]any
}

func (b uriBinding) Bind(r *http.Request, obj any) error {
	if err := binding(b.Params, obj); err != nil {
		return err
	}
	return validate(obj)
}
