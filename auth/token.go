package auth

import (
	"github.com/caixr9527/go-cloud/web"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

const TOKEN = "token"

type TokenAuth struct {
	UnAuthHandler func(context *web.Context)
	TokenName     string
	Key           []byte
}

func (ta *TokenAuth) Auth(context *web.Context) {
	if ta.TokenName == "" {
		ta.TokenName = TOKEN
	}
	token := context.R.Header.Get(ta.TokenName)
	if token == "" {
		token = context.Query(ta.TokenName)
		if token == "" {
			ta.unAuthHandler(context)
			return
		}
	}
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return ta.Key, nil
	})
	if err != nil {
		ta.UnAuthHandler(context)
		return
	}
	claims := t.Claims.(jwt.MapClaims)
	context.Set("claims", claims)
	context.Next()
}

func (ta *TokenAuth) unAuthHandler(context *web.Context) {
	if ta.UnAuthHandler != nil {
		ta.UnAuthHandler(context)
	} else {
		context.Fail(http.StatusUnauthorized, "Invalid token")
	}
	context.Abort()
}
