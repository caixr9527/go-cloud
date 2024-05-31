package log

import (
	"github.com/caixr9527/go-cloud/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"sync"
)

var Log *zap.Logger
var once sync.Once

func Init() {
	once.Do(func() {
		initLogger()
	})
}

func initLogger() {
	// 此处的配置是从我的项目配置文件读取的，读者可以根据自己的情况来设置
	//logPath := config.Cfg.Section("app").Key("logPath").String()
	//name := config.Cfg.Section("app").Key("name").String()
	//debug, err := config.Cfg.Section("app").Key("debug").Bool()
	//if err != nil {
	//	debug = false
	//}
	//logPath := "./logs/"
	loggerLevel := config.Cfg.Logger.Level
	if loggerLevel == "" {
		loggerLevel = "debug"
	}
	filename := config.Cfg.Logger.FileName
	if filename == "" {
		filename = "./logs/" + loggerLevel + ".log"
	}
	hook := lumberjack.Logger{
		Filename:   filename, // 日志文件路径
		MaxSize:    128,      // 每个日志文件保存的大小 单位:M
		MaxAge:     7,        // 文件最多保存多少天
		MaxBackups: 30,       // 日志文件最多保存多少个备份
		Compress:   true,     // 是否压缩
	}
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:       "msg",
		LevelKey:         "level",
		TimeKey:          "time",
		NameKey:          "logger",
		CallerKey:        "file",
		StacktraceKey:    "stacktrace",
		LineEnding:       zapcore.DefaultLineEnding,
		EncodeLevel:      zapcore.CapitalColorLevelEncoder,
		EncodeTime:       zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
		EncodeDuration:   zapcore.SecondsDurationEncoder,
		EncodeCaller:     zapcore.ShortCallerEncoder, // 短路径编码器
		EncodeName:       zapcore.FullNameEncoder,
		ConsoleSeparator: "|",
	}
	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(zap.DebugLevel)

	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
	cores := []zapcore.Core{
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(&hook), atomicLevel),
	}
	if loggerLevel == "debug" {
		atomicLevel.SetLevel(zap.DebugLevel)
		cores = append(cores, zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), atomicLevel))
	} else if loggerLevel == "info" {
		atomicLevel.SetLevel(zap.InfoLevel)
	} else if loggerLevel == "warn" {
		atomicLevel.SetLevel(zap.WarnLevel)
	} else if loggerLevel == "error" {
		atomicLevel.SetLevel(zap.ErrorLevel)
	}
	core := zapcore.NewTee(cores...)

	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()
	// 开启文件及行号
	development := zap.Development()
	// 构造日志
	Log = zap.New(core, caller, development)
	//Log = zap.New(core, caller)
}
