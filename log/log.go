package log

type LoggerLevel int

const (
	DEBUG LoggerLevel = iota
	INFO
	WARN
	ERROR
)
