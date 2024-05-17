package binding

import (
	"errors"
	"github.com/caixr9527/go-cloud/web/validator"
	jsoniter "github.com/json-iterator/go"
	"github.com/json-iterator/go/extra"
	"net/http"
	"reflect"
	"strings"
)

const defaultMaxMemory = 32 << 20

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
	FORM_MULTIPART = formMultipartBinding{}
	FORM_POST      = formPostBinding{}
	HEADER         = headerBinding{}
	MSG_PACK       = msgPackBinding{}
	PROTOBUF       = protobufBinding{}
	QUERY          = QueryBinding{}
	URI            = uriBinding{}
	YAML           = yamlBinding{}
	PLAIN          = plainBinding{}
)

func checkBody(body any) error {
	if body == nil {
		return errors.New("invalid request")
	}
	return nil
}

func validate(obj any) error {
	val := reflect.ValueOf(obj)
	for val.Kind() == reflect.Ptr && !val.IsNil() {
		val = val.Elem()
	}
	if val.Kind() == reflect.Struct {
		return validator.New().Struct(obj)
	}
	return nil
}
func parseMultipartForm(r *http.Request) error {
	if err := r.ParseMultipartForm(defaultMaxMemory); err != nil && !errors.Is(err, http.ErrNotMultipart) {
		return err
	}
	return nil
}

func mapping(datas map[string][]string, obj any) error {
	dataMap := make(map[string]any)
	for k, v := range datas {
		s := v[0]
		if strings.Contains(s, ",") {
			dataMap[k] = strings.Split(s, ",")
		} else if strings.Contains(s, ";") {
			dataMap[k] = strings.Split(s, ";")
		} else {
			dataMap[k] = s
		}
	}
	return binding(dataMap, obj)
}

func binding(dataMap map[string]any, obj any) error {
	extra.RegisterFuzzyDecoders()
	var marshal []byte
	var err error
	if marshal, err = jsoniter.Marshal(dataMap); err != nil {
		return err
	}
	if err = jsoniter.Unmarshal(marshal, obj); err != nil {
		return err
	}
	return err
}
