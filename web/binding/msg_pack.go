package binding

import "net/http"

type msgPackBinding struct {
}

func (b msgPackBinding) Bind(r *http.Request, obj any) error {

	return nil
}
