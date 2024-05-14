package binding

import "net/http"

type xmlBinding struct {
}

func (b xmlBinding) Bind(r *http.Request, obj any) error {

	return nil
}
