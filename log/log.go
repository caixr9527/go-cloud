package log

import (
	"fmt"
	"github.com/caixr9527/go-cloud/common/utils/stringUtils"
	"io"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

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
	Format(param *LoggingFormatParam) string
}

type LoggerWriter struct {
	Level LoggerLevel
	Out   io.Writer
}

func New() *Logger {
	return &Logger{}
}

func Default() *Logger {
	logger := New()
	logger.Level = DEBUG
	w := &LoggerWriter{
		Level: DEBUG,
		Out:   os.Stdout,
	}
	logger.Outs = append(logger.Outs, w)
	logger.Formatter = &TextFormatter{}
	return logger
}

func (l LoggerLevel) Level() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	default:
		return ""
	}
}

func (l *Logger) INFO(msg any) {
	l.Print(msg, INFO)
}

func (l *Logger) DEBUG(msg any) {
	l.Print(msg, DEBUG)
}

func (l *Logger) WARN(msg any) {
	l.Print(msg, WARN)
}

func (l *Logger) ERROR(msg any) {
	l.Print(msg, ERROR)
}
func (l *Logger) Print(msg any, level LoggerLevel) {
	if l.Level > level {
		return
	}

	param := &LoggingFormatParam{
		Level:        level,
		LoggerFields: l.LoggerFields,
		Msg:          msg,
	}
	str := l.Formatter.Format(param)
	for _, out := range l.Outs {
		if out.Out == os.Stdout {
			param.IsColor = true
			str = l.Formatter.Format(param)
			fmt.Fprintln(out.Out, str)
			l.checkFileSize(out)
		}
		if out.Level == -1 || level == out.Level {
			fmt.Fprintln(out.Out, str)
			l.checkFileSize(out)
		}
	}
}

func (l *Logger) WithFields(fields Fields) *Logger {
	return &Logger{
		Formatter:    l.Formatter,
		Outs:         l.Outs,
		Level:        l.Level,
		LoggerFields: fields,
	}
}

func (l *Logger) SetLogPath(logPath string) {
	l.logPath = logPath
	l.Outs = append(l.Outs, &LoggerWriter{
		Level: -1,
		Out:   FileWrite(path.Join(logPath, "all.log")),
	})
	l.Outs = append(l.Outs, &LoggerWriter{
		Level: DEBUG,
		Out:   FileWrite(path.Join(logPath, "debug.log")),
	})
	l.Outs = append(l.Outs, &LoggerWriter{
		Level: INFO,
		Out:   FileWrite(path.Join(logPath, "info.log")),
	})
	l.Outs = append(l.Outs, &LoggerWriter{
		Level: WARN,
		Out:   FileWrite(path.Join(logPath, "warn.log")),
	})
	l.Outs = append(l.Outs, &LoggerWriter{
		Level: ERROR,
		Out:   FileWrite(path.Join(logPath, "error.log")),
	})
}

func (l *Logger) checkFileSize(writer *LoggerWriter) {
	logFile := writer.Out.(*os.File)
	if logFile != nil {
		stat, err := logFile.Stat()
		if err != nil {
			log.Println(err)
			return
		}
		size := stat.Size()
		if l.LogFileSize <= 0 {
			l.LogFileSize = 100 << 20
		}
		if size >= l.LogFileSize {
			// todo 需要优化，应该一直往info.log文件里面写，满了再归档到另一个文件下
			// todo 可添加，按天归档
			_, name := path.Split(stat.Name())
			fileName := name[0:strings.Index(name, ".")]
			write := FileWrite(path.Join(l.logPath, stringUtils.JoinStrings(fileName, ".", time.Now().UnixMilli(), ".log")))
			writer.Out = write
		}
	}
}

func FileWrite(name string) io.Writer {
	w, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	return w
}
