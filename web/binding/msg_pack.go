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
	err := decoder.Decode(obj)
	if err != nil {
		return err
	}
	return validate(obj)
}
