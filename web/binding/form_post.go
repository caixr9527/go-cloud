package binding

import (
	"net/http"
)

type formPostBinding struct {
}

func (b formPostBinding) Bind(r *http.Request, obj any) error {

	return validate(obj)
}
