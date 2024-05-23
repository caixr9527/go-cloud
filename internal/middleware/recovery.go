package middleware

import (
	"fmt"
	logger "github.com/caixr9527/go-cloud/log"
	"github.com/caixr9527/go-cloud/web"
	"net/http"
	"runtime"
	"strings"
)

func Recovery(context *web.Context) {
	defer func() {
		if err := recover(); err != nil {
			logger.Log.Error(detailMsg(err))
			context.Fail(http.StatusInternalServerError, "Internal Server Error")
		}
	}()
	context.Next()
}

func detailMsg(err any) string {
	var pcs [32]uintptr
	n := runtime.Callers(0, pcs[:])
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%v", err))
	for _, pc := range pcs[0:n] {
		fn := runtime.FuncForPC(pc)
		file, l := fn.FileLine(pc)
		sb.WriteString(fmt.Sprintf("\n\t%s:%d", file, l))

	}
	return sb.String()
}
