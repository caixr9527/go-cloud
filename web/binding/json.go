package binding

import "net/http"

type jsonBinding struct {
	DisallowUnknownFields bool
	IsValidate            bool
}

func (b jsonBinding) Bind(r *http.Request, obj any) error {

	return nil
}
