package log

type TextFormatter struct {
}

func (f *TextFormatter) Format(param *LoggingFormatParam) string {
	return ""
}
