package binding

import (
	"net/http"
)

type headerBinding struct {
}

func (b headerBinding) Bind(r *http.Request, obj any) error {
	if err := mapping(r.Header, obj); err != nil {
		return err
	}
	return validate(obj)
}
