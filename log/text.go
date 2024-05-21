package log

import (
	"fmt"
	"strings"
	"time"
)

type TextFormatter struct {
}

func (f *TextFormatter) Format(param *LoggingFormatParam) string {
	now := time.Now()
	fieldsString := ""
	if param.LoggerFields != nil {
		var sb strings.Builder
		for key, val := range param.LoggerFields {
			fmt.Fprintf(&sb, "%s=%v%s", key, val, ",")
		}
		fieldsString = sb.String()[0 : sb.Len()-1]
	}
	var msgInfo = "\n"
	if param.Level == ERROR {
		msgInfo = "\n Error Cause By :"
	}
	if param.IsColor {
		levelColor := f.LevelColor(param.Level)
		msgColor := f.MsgColor(param.Level)
		return fmt.Sprintf("[go-cloud] | %s [%s] %s | %v  | %s %s%v %s ｜ %s",
			levelColor, param.Level.Level(), Reset,
			now.Format("2006/01/02 - 15:04:05"),
			msgColor, msgInfo, param.Msg, Reset,
			fieldsString)
	}
	return fmt.Sprintf("[go-cloud] | [%s] | %v  | %s %v | %s",
		param.Level.Level(),
		now.Format("2006/01/02 - 15:04:05"),
		msgInfo, param.Msg, fieldsString)
}

func (f *TextFormatter) LevelColor(level LoggerLevel) interface{} {
	switch level {
	case DEBUG:
		return Blue
	case INFO:
		return Green
	case WARN:
		return Yellow
	case ERROR:
		return Red
	default:
		return ""
	}
}

func (f *TextFormatter) MsgColor(level LoggerLevel) interface{} {
	switch level {
	case DEBUG:
		return Blue
	case INFO:
		return Green
	case WARN:
		return Yellow
	case ERROR:
		return Red
	default:
		return ""
	}
}
