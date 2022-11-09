package logger

import "io"

type LoggerOption interface {
	Apply(l *Logger)
}

type LoggerOptionFunc func(l *Logger)

func (lof LoggerOptionFunc) Apply(l *Logger) {
	lof(l)
}

// WithNowFunc overrides default function used to determine current time.
// Intended to be used in tests only.
func WithNowFunc(nowFunc NowFunc) LoggerOption {
	return LoggerOptionFunc(func(l *Logger) {
		l.now = nowFunc
	})
}

// WithOutput overrides default output the logs are written to.
func WithOutput(output io.Writer) LoggerOption {
	return LoggerOptionFunc(func(l *Logger) {
		l.output = output
	})
}

// WithLevel sets minimum level for filtering logs
func WithLevel(level Level) LoggerOption {
	return LoggerOptionFunc(func(l *Logger) {
		l.level = level
	})
}

// WithReportCaller allows enabling/disabling including calling method in the log entry
func WithReportCaller(enable bool) LoggerOption {
	return LoggerOptionFunc(func(l *Logger) {
		l.reportCaller = enable
	})
}
