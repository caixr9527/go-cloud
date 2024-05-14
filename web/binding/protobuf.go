package binding

import "net/http"

type protobufBinding struct {
}

func (b protobufBinding) Bind(r *http.Request, obj any) error {

	return nil
}
