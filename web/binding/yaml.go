package binding

import "net/http"

type yamlBinding struct {
}

func (b yamlBinding) Bind(r *http.Request, obj any) error {

	return nil
}
