package binding

import (
	"io"
	"net/http"
	"net/url"
)

type formPostBinding struct {
}

func (b formPostBinding) Bind(r *http.Request, obj any) error {
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return err
	}
	datas, err := url.ParseQuery(string(body))
	if err != nil {
		return err
	}
	if err = mapping(datas, obj); err != nil {
		return err
	}
	return validate(obj)
}
