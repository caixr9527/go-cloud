package binding

import (
	"encoding/xml"
	"net/http"
)

type xmlBinding struct {
}

func (b xmlBinding) Bind(r *http.Request, obj any) error {
	body := r.Body
	if err := checkBody(body); err != nil {
		return err
	}
	decoder := xml.NewDecoder(body)
	err := decoder.Decode(obj)
	if err != nil {
		return err
	}
	return validate(obj)
}
