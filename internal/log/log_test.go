package log_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/coopnorge/go-datadog-lib/v2/internal/log"
	coopLogger "github.com/coopnorge/go-logger"

	"github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetup(t *testing.T) {
	logger, err := log.NewLogger(log.WithGlobalLogger())
	require.NoError(t, err)
	tracer.UseLogger(logger)
}

func TestGlobalLogger(t *testing.T) {
	output := &strings.Builder{}
	coopLogger.ConfigureGlobalLogger(coopLogger.WithLevel(coopLogger.LevelDebug), coopLogger.WithOutput(output))
	logger, err := log.NewLogger(log.WithGlobalLogger())
	require.NoError(t, err)

	tests := []struct {
		level string
		msg   string
	}{
		{"error", "Datadog Tracer v1.63.0 ERROR This is a test"},
		{"warning", "Datadog Tracer v1.63.0 WARN This is a test"},
		{"info", "Datadog Tracer v1.63.0 INFO This is a test"},
		{"debug", "Datadog Tracer v1.63.0 DEBUG This is a test"},
		{"warning", "Datadog Tracer v1.63.0 This is a test for fallback level"},
	}
	for _, test := range tests {
		t.Run(test.level, func(t *testing.T) {
			output.Reset()
			logger.Log(test.msg)
			assert.Contains(t, output.String(), fmt.Sprintf("\"level\":\"%v\"", test.level))
			assert.Contains(t, output.String(), test.msg)
		})
	}
}

func TestCustomLogger(t *testing.T) {
	output := &strings.Builder{}
	logger, err := log.NewLogger(log.WithLogger(coopLogger.New(coopLogger.WithLevel(coopLogger.LevelDebug), coopLogger.WithOutput(output))))
	require.NoError(t, err)

	tests := []struct {
		level string
		msg   string
	}{
		{"error", "Datadog Tracer v1.63.0 ERROR This is a test"},
		{"warning", "Datadog Tracer v1.63.0 WARN This is a test"},
		{"info", "Datadog Tracer v1.63.0 INFO This is a test"},
		{"debug", "Datadog Tracer v1.63.0 DEBUG This is a test"},
		{"warning", "Datadog Tracer v1.63.0 This is a test for fallback level"},
	}
	for _, test := range tests {
		t.Run(test.level, func(t *testing.T) {
			output.Reset()
			logger.Log(test.msg)
			assert.Contains(t, output.String(), fmt.Sprintf("\"level\":\"%v\"", test.level))
			assert.Contains(t, output.String(), test.msg)
		})
	}
}
