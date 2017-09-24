package logger

type Logger interface {
	Log(location string, message string)
}
