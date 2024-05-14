package binding

import (
	"github.com/caixr9527/go-cloud/web/validator"
	"net/http"
)

type msgPackBinding struct {
}

func (b msgPackBinding) Bind(r *http.Request, obj any) error {

	return validator.New().Struct(obj)
}
