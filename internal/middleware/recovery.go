package middleware

import (
	"fmt"
	"github.com/caixr9527/go-cloud/component/factory"
	"github.com/caixr9527/go-cloud/log"
	"github.com/caixr9527/go-cloud/web"
	"net/http"
	"runtime"
	"strings"
)

func Recovery(context *web.Context) {
	defer func() {
		logger := factory.Get(&log.Log{})
		if err := recover(); err != nil {
			logger.Error(detailMsg(err))
			context.Fail(http.StatusInternalServerError, "Internal Server Error")
			context.Abort()
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
