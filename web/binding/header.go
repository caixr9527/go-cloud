package binding

import "net/http"

type headerBinding struct {
}

func (b headerBinding) Bind(r *http.Request, obj any) error {

	return nil
}
