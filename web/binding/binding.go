package binding

import "net/http"

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
