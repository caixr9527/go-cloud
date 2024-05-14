package binding

import "net/http"

type QueryBinding struct {
}

func (b QueryBinding) Bind(r *http.Request, obj any) error {

	return nil
}
