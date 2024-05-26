package auth

import (
	"github.com/caixr9527/go-cloud/web"
	"net/http"
)

type BasicAuth struct {
	UnAuthHandler func(context *web.Context)
	AuthKeys      map[string]any
	Realm         string
}

func (b *BasicAuth) Auth(context *web.Context) {
	username, password, ok := context.R.BasicAuth()
	if !ok {
		b.unAuthHandler(context)
		return
	}
	pwd, exist := b.AuthKeys[username]
	if !exist {
		b.unAuthHandler(context)
		return
	}
	if pwd != password {
		b.unAuthHandler(context)
		return
	}
	context.Set("user", username)
	context.Set("pwd", pwd)
	context.Next()
}

func (b *BasicAuth) unAuthHandler(context *web.Context) {
	if b.UnAuthHandler != nil {
		b.UnAuthHandler(context)
	} else {
		context.W.Header().Set("WWW-Authenticate", b.Realm)
		context.Fail(http.StatusUnauthorized, "Authentication failed")
	}
	context.Abort()
}
