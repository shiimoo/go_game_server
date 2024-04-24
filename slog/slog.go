package slog

var defaultLogger = NewLogger()

func Default() *Logger {
	return defaultLogger
}

func Log(v ...any) {
	Default().Log(v...)
}
