package binding

import (
	"net/http"
)

type formBinding struct {
}

func (b formBinding) Bind(r *http.Request, obj any) error {

	return validate(obj)
}
