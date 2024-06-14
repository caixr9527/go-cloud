package middleware

import (
	"fmt"
	"github.com/caixr9527/go-cloud/component/factory"
	"github.com/caixr9527/go-cloud/web"
	"go.uber.org/zap"
	"net"
	"strings"
	"time"
)

func Logging(context *web.Context) {
	r := context.R
	start := time.Now()
	path := r.URL.Path
	raw := r.URL.RawQuery

	context.Next()

	stop := time.Now()
	stop.Sub(start)
	latency := stop.Sub(start)
	ip, _, _ := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	clientIp := net.ParseIP(ip)
	method := r.Method
	statusCode := context.StatusCode
	if raw != "" {
		path = path + "?" + raw
	}
	logger := factory.Get(&zap.Logger{})
	logger.Debug(fmt.Sprintf("ip: %s, method: %s, path: %s, status: %3d, cost: %v ", clientIp, method, path, statusCode, latency))
}
