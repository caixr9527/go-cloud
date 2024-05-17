package binding

import (
	"fmt"
	"io"
	"net/http"
	"reflect"
)

type plainBinding struct {
}

func (b plainBinding) Bind(r *http.Request, obj any) error {
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return err
	}
	valueOf := reflect.ValueOf(obj)
	for valueOf.Kind() == reflect.Ptr {
		if valueOf.IsNil() {
			return nil
		}
		valueOf = valueOf.Elem()
	}
	if valueOf.Kind() == reflect.String {
		valueOf.SetString(string(body))
		return nil
	}

	if _, ok := valueOf.Interface().([]byte); ok {
		valueOf.SetBytes(body)
		return nil
	}

	return fmt.Errorf("type (%T) unknown type", valueOf)
}
