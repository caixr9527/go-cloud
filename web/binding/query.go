package binding

import (
	"net/http"
)

type QueryBinding struct {
}

func (b QueryBinding) Bind(r *http.Request, obj any) error {
	values := r.URL.Query()
	if err := mapping(values, obj); err != nil {
		return err
	}
	return validate(obj)
}
