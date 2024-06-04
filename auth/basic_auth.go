package auth

import (
	"github.com/caixr9527/go-cloud/config"
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
	if config.Cfg.BasicAuth.Username != username || password != config.Cfg.BasicAuth.Password {
		unAuth(context, "Authentication failed")
		return
	}
	context.Set("user", username)
	context.Set("pwd", password)
}

func unAuth(context *web.Context, msg string) {
	context.W.Header().Set("WWW-Authenticate", config.Cfg.BasicAuth.Realm)
	context.Fail(http.StatusUnauthorized, msg)
	context.Abort()
}
