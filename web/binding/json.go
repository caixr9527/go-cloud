package binding

import (
	"encoding/json"
	"net/http"
)

type jsonBinding struct {
	DisallowUnknownFields bool
}

func (b jsonBinding) Bind(r *http.Request, obj any) error {
	body := r.Body
	if err := checkBody(body); err != nil {
		return err
	}
	decoder := json.NewDecoder(body)
	if b.DisallowUnknownFields {
		decoder.DisallowUnknownFields()
	}
	err := decoder.Decode(obj)
	if err != nil {
		return err
	}
	return validate(obj)
}
