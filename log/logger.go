package log

import (
	"github.com/caixr9527/go-cloud/common/utils"
	"github.com/caixr9527/go-cloud/config"
	"github.com/caixr9527/go-cloud/factory"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"math"
	"os"
	"sync"
)

var once sync.Once

type Log struct {
	*zap.Logger
}

func (l *Log) Create() {
	once.Do(func() {
		l.initLogger()
	})
}

func init() {
	factory.RegisterComponent(&Log{})
}
func (l *Log) Refresh() {
	l.initLogger()
}
func (l *Log) Order() int {
	return math.MinInt + 2
}

func (l *Log) Destroy() {

}

func (l *Log) Name() string {
	return utils.ObjName(l)
}

func (l *Log) initLogger() {
	configuration := factory.Get(&config.Configuration{})
	loggerLevel := configuration.Logger.Level
	if loggerLevel == "" {
		loggerLevel = "debug"
	}
	filename := configuration.Logger.FileName
	if filename == "" {
		filename = "./logs/" + loggerLevel + ".log"
	}
	maxAge := configuration.Logger.MaxAge
	if maxAge == 0 {
		maxAge = 7
	}
	maxSize := configuration.Logger.MaxSize
	if maxSize == 0 {
		maxSize = 100
	}
	maxBackups := configuration.Logger.MaxBackups
	if maxBackups == 0 {
		maxBackups = 7
	}
	hook := lumberjack.Logger{
		Filename:   filename,                      // 日志文件路径
		MaxSize:    int(maxSize),                  // 每个日志文件保存的大小 单位:M
		MaxAge:     int(maxAge),                   // 文件最多保存多少天
		MaxBackups: int(maxBackups),               // 日志文件最多保存多少个备份
		Compress:   configuration.Logger.Compress, // 是否压缩
	}
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:       "msg",
		LevelKey:         "level",
		TimeKey:          "time",
		NameKey:          "Log",
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
	active := configuration.Cloud.Active
	//var Log *zap.Logger
	if active == config.DEV {
		// 开启文件及行号
		development := zap.Development()
		// 构造日志
		l.Logger = zap.New(core, caller, development)
	} else {
		l.Logger = zap.New(core, caller)
	}
	factory.Create(l)
}
