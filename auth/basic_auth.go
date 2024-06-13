package auth

import (
	"github.com/caixr9527/go-cloud/component/factory"
	"github.com/caixr9527/go-cloud/web"
	"net/http"
)

// YWRtaW46Z29fY2xvdWQ=

func BasicAuth(context *web.Context) {
	username, password, ok := context.R.BasicAuth()
	if !ok {
		unAuth(context, "basic auth require")
		return
	}
	if factory.GetConf().BasicAuth.Username != username || password != factory.GetConf().BasicAuth.Password {
		unAuth(context, "Authentication failed")
		return
	}
	context.Set("user", username)
	context.Set("pwd", password)
}

func unAuth(context *web.Context, msg string) {
	context.W.Header().Set("WWW-Authenticate", factory.GetConf().BasicAuth.Realm)
	context.Fail(http.StatusUnauthorized, msg)
	context.Abort()
}
