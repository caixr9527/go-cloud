package binding

import (
	"encoding/json"
	"errors"
	"net/http"
)

const defaultMaxMemory = 32 << 20

type formMultipartBinding struct {
}

func (b formMultipartBinding) Bind(r *http.Request, obj any) error {
	if err := r.ParseMultipartForm(defaultMaxMemory); err != nil && !errors.Is(err, http.ErrNotMultipart) {
		return err
	}
	form := r.Form
	var marshal []byte
	var err error
	if marshal, err = json.Marshal(form); err != nil {
		return err
	}
	if err = json.Unmarshal(marshal, obj); err != nil {
		return err
	}
	return validate(obj)
}
