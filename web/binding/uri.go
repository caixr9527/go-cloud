package binding

import "net/http"

type uriBinding struct {
}

func (b uriBinding) Bind(r *http.Request, obj any) error {

	return nil
}
