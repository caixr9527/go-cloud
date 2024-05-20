package log

import "io"

type LoggerLevel int

var (
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Reset  = "\033[0m"
)

const (
	DEBUG LoggerLevel = iota
	INFO
	WARN
	ERROR
)

type Logger struct {
	Formatter    LoggingFormatter
	Level        LoggerLevel
	Outs         []*LoggerWriter
	LoggerFields Fields
	logPath      string
	LogFileSize  int64
}

type Fields map[string]any
type LoggingFormatParam struct {
	Level        LoggerLevel
	IsColor      bool
	LoggerFields Fields
	Msg          any
}

type LoggingFormatter interface {
	Format(param *LoggingFormatParam)
}

type LoggerWriter struct {
	Level LoggerLevel
	Out   io.Writer
}

func New() *Logger {
	return &Logger{}
}
