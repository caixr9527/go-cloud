package log

import (
	"fmt"
	"github.com/caixr9527/go-cloud/web"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"net"
	"os"
	"strings"
	"time"
)

var Log *logger

type logger struct {
	log *zap.Logger
}

func Info(msg string, fields ...zap.Field) {
	defer Log.log.Sync()
	Log.log.Info(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	defer Log.log.Sync()
	Log.log.Error(msg, fields...)
}

func Debug(msg string, fields ...zap.Field) {
	defer Log.log.Sync()
	Log.log.Debug(msg, fields...)
}
func Warn(msg string, fields ...zap.Field) {
	defer Log.log.Sync()
	Log.log.Warn(msg, fields...)
}

func init() {
	InitLogger()
}

func InitLogger() {
	// 此处的配置是从我的项目配置文件读取的，读者可以根据自己的情况来设置
	//logPath := config.Cfg.Section("app").Key("logPath").String()
	//name := config.Cfg.Section("app").Key("name").String()
	//debug, err := config.Cfg.Section("app").Key("debug").Bool()
	//if err != nil {
	//	debug = false
	//}
	logPath := "./logs/"
	//name := "logTest"
	debug := true
	infoHook := lumberjack.Logger{
		Filename:   logPath + "info.log", // 日志文件路径
		MaxSize:    128,                  // 每个日志文件保存的大小 单位:M
		MaxAge:     7,                    // 文件最多保存多少天
		MaxBackups: 30,                   // 日志文件最多保存多少个备份
		Compress:   true,                 // 是否压缩
	}
	debugHook := lumberjack.Logger{
		Filename:   logPath + "debug.log", // 日志文件路径
		MaxSize:    128,                   // 每个日志文件保存的大小 单位:M
		MaxAge:     7,                     // 文件最多保存多少天
		MaxBackups: 30,                    // 日志文件最多保存多少个备份
		Compress:   true,                  // 是否压缩
	}
	warnHook := lumberjack.Logger{
		Filename:   logPath + "warn.log", // 日志文件路径
		MaxSize:    128,                  // 每个日志文件保存的大小 单位:M
		MaxAge:     7,                    // 文件最多保存多少天
		MaxBackups: 30,                   // 日志文件最多保存多少个备份
		Compress:   true,                 // 是否压缩
	}
	errorHook := lumberjack.Logger{
		Filename:   logPath + "error.log", // 日志文件路径
		MaxSize:    128,                   // 每个日志文件保存的大小 单位:M
		MaxAge:     7,                     // 文件最多保存多少天
		MaxBackups: 30,                    // 日志文件最多保存多少个备份
		Compress:   true,                  // 是否压缩
	}
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "file",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder, // 短路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
	cores := make([]zapcore.Core, 0)
	cores = append(cores, zapcore.NewCore(consoleEncoder, zapcore.AddSync(&infoHook), zap.InfoLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(&errorHook), zap.ErrorLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(&debugHook), zap.DebugLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(&warnHook), zap.WarnLevel))
	if debug {
		cores = append(cores, zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zap.DebugLevel))
	}
	core := zapcore.NewTee(cores...)

	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()
	// 开启文件及行号
	development := zap.Development()
	// 构造日志
	Log = &logger{log: zap.New(core, caller, development)}
}

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
	Debug(fmt.Sprintf("ip: %s, method: %s, path: %s, status: %3d, cost: %v ", clientIp, method, path, statusCode, latency))
}
