package binding

import (
	"errors"
	"github.com/caixr9527/go-cloud/web/validator"
	"net/http"
	"reflect"
)

type Binding interface {
	Bind(r *http.Request, obj any) error
}

const (
	MIMEJSON              = "application/json"
	MIMEHTML              = "text/html"
	MIMEXML               = "application/xml"
	MIMEXML2              = "text/xml"
	MIMEPLAIN             = "text/plain"
	MIMEPOSTFORM          = "application/x-www-form-urlencoded"
	MIMEMultipartPOSTForm = "multipart/form-data"
	MIMEPROTOBUF          = "application/x-protobuf"
	MIMEMSGPACK           = "application/x-msgpack"
	MIMEMSGPACK2          = "application/msgpack"
	MIMEYAML              = "application/x-yaml"
)

var (
	JSON           = jsonBinding{}
	XML            = xmlBinding{}
	FORM           = formBinding{}
	FORM_MULTIPART = formMultipartBinding{}
	FORM_POST      = formPostBinding{}
	HEADER         = headerBinding{}
	MSG_PACK       = msgPackBinding{}
	PROTOBUF       = protobufBinding{}
	QUERY          = QueryBinding{}
	URI            = uriBinding{}
	YAML           = yamlBinding{}
)

func checkBody(body any) error {
	if body == nil {
		return errors.New("invalid request")
	}
	return nil
}

func validate(obj any) error {
	typeOf := reflect.TypeOf(obj)
	if typeOf.Kind() == reflect.Struct {
		return validator.New().Struct(obj)
	}
	return nil
}
