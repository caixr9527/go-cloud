package auth

import "github.com/caixr9527/go-cloud/web"

type Authentication interface {
	Auth(context *web.Context)
}
