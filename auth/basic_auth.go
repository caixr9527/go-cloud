package auth

import (
	"github.com/caixr9527/go-cloud/component/factory"
	"github.com/caixr9527/go-cloud/config"
	"github.com/caixr9527/go-cloud/web"
	"net/http"
)

// YWRtaW46Z29fY2xvdWQ=

func BasicAuth(context *web.Context) {
	configuration := factory.Get(config.Configuration{})
	username, password, ok := context.R.BasicAuth()
	if !ok {
		unAuth(context, "basic auth require")
		return
	}
	if configuration.BasicAuth.Username != username || password != configuration.BasicAuth.Password {
		unAuth(context, "Authentication failed")
		return
	}
	context.Set("user", username)
	context.Set("pwd", password)
}

func unAuth(context *web.Context, msg string) {
	configuration := factory.Get(config.Configuration{})
	context.W.Header().Set("WWW-Authenticate", configuration.BasicAuth.Realm)
	context.Fail(http.StatusUnauthorized, msg)
	context.Abort()
}
