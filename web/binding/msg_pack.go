package binding

import (
	"github.com/vmihailenco/msgpack/v5"
	"net/http"
)

type msgPackBinding struct {
}

func (b msgPackBinding) Bind(r *http.Request, obj any) error {
	body := r.Body
	if err := checkBody(body); err != nil {
		return err
	}
	decoder := msgpack.NewDecoder(body)
	if err := decoder.Decode(obj); err != nil {
		return err
	}
	return validate(obj)
}
