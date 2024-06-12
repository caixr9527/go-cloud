package cors

import (
	"github.com/caixr9527/go-cloud/web"
	"net/http"
)

func Cors(context *web.Context) {
	if context.R.Method != "OPTIONS" && context.R.Method != "options" {
		return
	} else {
		context.W.Header().Set("Access-Control-Allow-Origin", "*")
		context.W.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		context.W.Header().Set("Access-Control-Allow-Headers", "authorization, origin, content-type, accept")
		context.W.Header().Set("Allow", "HEAD,GET,POST,PUT,PATCH,DELETE,OPTIONS")
		context.W.Header().Set("Content-Type", "application/json")
		context.JSON(http.StatusOK)
		context.Abort()
	}
}
