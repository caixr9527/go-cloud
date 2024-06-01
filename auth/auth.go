package auth

import (
	"github.com/caixr9527/go-cloud/config"
	"github.com/caixr9527/go-cloud/web"
	"time"
)

type Authentication interface {
	Auth(context *web.Context)
}

var (
	// YWRtaW46Z29fY2xvdWQ=
	Basic = &BasicAuth{AuthKeys: map[string]any{config.Cfg.BasicAuth.Username: config.Cfg.BasicAuth.Password}}
	Token = &TokenAuth{
		TokenName: "token",
		JwtConfig: JwtConfig{
			Alg:          config.Cfg.Jwt.Alg,
			TokenTimeout: config.Cfg.Jwt.TokenTimeout * time.Second,
			RefreshKey:   []byte(config.Cfg.Jwt.RefreshKey),
			Key:          []byte(config.Cfg.Jwt.SecretKey),
			Whitelist:    config.Cfg.Jwt.Whitelist,
		},
	}
)
