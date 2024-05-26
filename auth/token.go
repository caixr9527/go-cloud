package auth

import (
	"github.com/caixr9527/go-cloud/web"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

const TOKEN = "token"

type TokenAuth struct {
	JwtConfig
	UnAuthHandler func(context *web.Context)
	TokenName     string
}

type JwtResponse struct {
	Token        string
	RefreshToken string
}

type JwtConfig struct {
	Alg                 string
	TokenTimeout        time.Duration
	RefreshTokenTimeout time.Duration
	RefreshKey          []byte
	Key                 []byte
	Whitelist           []string
}

type JwtToken struct {
	JwtConfig
}

func (jt *JwtToken) CreateToken(claims map[string]any) (*JwtResponse, error) {
	signingMethod := jwt.GetSigningMethod(jt.Alg)
	token := jwt.New(signingMethod)
	mapClaims := token.Claims.(jwt.MapClaims)
	for k, v := range claims {
		mapClaims[k] = v
	}

	claims["exp"] = time.Now().Add(jt.TokenTimeout).Unix()
	claims["iat"] = time.Now().Unix()
	tokenString, err := token.SignedString(jt.Key)
	if err != nil {
		return nil, err
	}
	token.Claims = jwt.MapClaims{
		"exp": time.Now().Add(jt.RefreshTokenTimeout).Unix(),
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
