package logger

import "github.com/sirupsen/logrus"

type Level uint8

const (
	// LevelFatal is to be used to log predictable errors that make the service unsuable, eg misconfiguration. After logging, the app will be shut down.
	LevelFatal = iota
	// LevelError  isto be used for recoverable errors that limit the service's functionality, eg timeouts.
	LevelError
	// LevelWarn is to be used for non-critical errors that may require some attention.
	LevelWarn
	// LevelInfo is to be used for monitoring successful interactions, eg run time or job result.
	LevelInfo
	// LevelDebug should only be used in dev/test environments.
	LevelDebug
)

func mapLevelToLogrusLevel(l Level) logrus.Level {
	switch l {
	case LevelFatal:
		return logrus.FatalLevel
	case LevelError:
		return logrus.ErrorLevel
	case LevelWarn:
		return logrus.WarnLevel
	case LevelInfo:
		return logrus.InfoLevel
	case LevelDebug:
		return logrus.DebugLevel
	}
	// should never get here
	return logrus.DebugLevel
}
