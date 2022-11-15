package logger

import "io"

var globalLogger = New()

func ConfigureGlobalLogger(opts ...LoggerOption) {
	globalLogger.applyOptions(opts...)
}

func Global() *Logger {
	return globalLogger
}

// WithFields creates log entry using global logger
func WithFields(fields Fields) Entry {
	return globalLogger.WithFields(fields)
}

// Infof uses global logger to log payload on "info" level
func Info(args ...interface{}) {
	globalLogger.Info(args...)
}

// Infof uses global logger to log payload on "info" level
func Infof(format string, args ...interface{}) {
	globalLogger.Infof(format, args...)
}

// Error uses global logger to log payload on "error" level
func Error(args ...interface{}) {
	globalLogger.Error(args...)
}

// Errorf uses global logger to log payload on "error" level
func Errorf(format string, args ...interface{}) {
	globalLogger.Errorf(format, args...)
}

// Debug uses global logger to log payload on "debug" level
func Debug(args ...interface{}) {
	globalLogger.Debug(args...)
}

// Debugf uses global logger to log payload on "debug" level
func Debugf(format string, args ...interface{}) {
	globalLogger.Debugf(format, args...)
}

// Warn uses global logger to log payload on "warn" level
func Warn(args ...interface{}) {
	globalLogger.Warn(args...)
}

// Warnf uses global logger to log payload on "warn" level
func Warnf(format string, args ...interface{}) {
	globalLogger.Warnf(format, args...)
}

// Fatal uses global logger to log payload on "fatal" level
func Fatal(args ...interface{}) {
	globalLogger.Fatal(args...)
}

// Fatalf uses global logger to log payload on "fatal" level
func Fatalf(format string, args ...interface{}) {
	globalLogger.Fatalf(format, args...)
}

// SetNowFunc sets `now` func user by global logger
func SetNowFunc(nowFunc NowFunc) {
	ConfigureGlobalLogger(WithNowFunc(nowFunc))
}

// SetOutput changes global logger's output
func SetOutput(output io.Writer) {
	ConfigureGlobalLogger(WithOutput(output))
}

// SetLevel sets minimum log level for global logger
func SetLevel(level Level) {
	ConfigureGlobalLogger(WithLevel(level))
}

// SetReportCaller allows controling if caller info should be attached to logs by global logger
func SetReportCaller(enable bool) {
	ConfigureGlobalLogger(WithReportCaller(enable))
}
