package auth

import (
	"github.com/caixr9527/go-cloud/config"
	"github.com/caixr9527/go-cloud/web"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

const TOKEN = "token"

type JwtResponse struct {
	Token        string
	RefreshToken string
}

type JwtToken struct {
	Alg          string
	TokenTimeout time.Duration
	RefreshKey   []byte
	Key          []byte
}

func (jt *JwtToken) CreateToken(claims map[string]any) (*JwtResponse, error) {
	if jt.Alg == "" {
		jt.Alg = "HS256"
	}
	signingMethod := jwt.GetSigningMethod(jt.Alg)
	token := jwt.New(signingMethod)
	mapClaims := token.Claims.(jwt.MapClaims)
	for k, v := range claims {
		mapClaims[k] = v
	}

	mapClaims["exp"] = time.Now().Add(jt.TokenTimeout).Unix()
	mapClaims["iat"] = time.Now().Unix()
	tokenString, err := token.SignedString(jt.Key)
	if err != nil {
		return nil, err
	}
	token.Claims = jwt.MapClaims{
		"exp": time.Now().Add(2 * jt.TokenTimeout).Unix(),
	}
	refreshToken, err := token.SignedString(jt.RefreshKey)
	if err != nil {
		return nil, err
	}
	return &JwtResponse{
		Token:        tokenString,
		RefreshToken: refreshToken,
	}, nil
}

func Token(context *web.Context) {
	//if sliceUtils.Contains(config.Cfg.Jwt.Whitelist, context.R.URL.Path) {
	//	return
	//}
	if config.Cfg.Jwt.Header == "" {
		config.Cfg.Jwt.Header = TOKEN
	}
	token := context.R.Header.Get(config.Cfg.Jwt.Header)
	if token == "" {
		token = context.Query(config.Cfg.Jwt.Header)
		if token == "" {
			unAuthHandler(context, "token require")
			return
		}
	}
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Cfg.Jwt.SecretKey), nil
	})
	if err != nil {
		unAuthHandler(context, err.Error())
		return
	}
	claims := t.Claims.(jwt.MapClaims)
	context.Set("claims", claims)
}

func unAuthHandler(context *web.Context, msg string) {
	context.Fail(http.StatusUnauthorized, msg)
	context.Abort()
}
