package binding

import (
	"net/http"
)

type formMultipartBinding struct {
}

func (b formMultipartBinding) Bind(r *http.Request, obj any) error {
	var err error
	if err = parseMultipartForm(r); err != nil {
		return err
	}
	form := r.Form
	if err = mapping(form, obj); err != nil {
		return err
	}
	return validate(obj)
}
