package auth

import "github.com/caixr9527/go-cloud/web"

type Authentication interface {
	Auth(context *web.Context)
}

var (
	// YWRtaW46Z29fY2xvdWQ=
	Basic = &BasicAuth{AuthKeys: map[string]any{"admin": "go_cloud"}}
	Token = &TokenAuth{
		TokenName: "token",
		JwtConfig: JwtConfig{
			Key: []byte("go_cloud"),
		},
	}
)
