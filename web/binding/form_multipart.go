package binding

import "net/http"

type formMultipartBinding struct {
}

func (b formMultipartBinding) Bind(r *http.Request, obj any) error {

	return nil
}
