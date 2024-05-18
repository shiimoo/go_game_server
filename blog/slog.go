package blog

var defaultLogger = NewLogger("defaultLogger")

func Default() *Logger {
	return defaultLogger
}

func Log(v ...any) {
	Default().Log(v...)
}

func Logf(format string, v ...any) {
	Default().Logf(format, v...)
}

func Debug(v ...any) {
	Default().Debug(v...)
}

func Debugf(format string, v ...any) {
	Default().Debugf(format, v...)
}

func Info(v ...any) {
	Default().Info(v...)
}

func Infof(format string, v ...any) {
	Default().Infof(format, v...)
}

func Warn(v ...any) {
	Default().Warn(v...)
}

func Warnf(format string, v ...any) {
	Default().Warnf(format, v...)
}

func Error(v ...any) {
	Default().Error(v...)
}

func Errorf(format string, v ...any) {
	Default().Errorf(format, v...)
}

func Fatal(v ...any) {
	Default().Fatal(v...)
}

func Fatalf(format string, v ...any) {
	Default().Fatalf(format, v...)
}
