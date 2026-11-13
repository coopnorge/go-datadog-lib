// Package errors reports errors to Datadog spans.
//
// Deprecated: The whole module is deprecated. Use span.RecordError and
// span.SetStatus from go.opentelemetry.io/otel/trace and
// go.opentelemetry.io/otel/codes. See the root package documentation for the
// migration guide and ADR.
package errors //nolint:revive

// ErrorHandler allows for handling of error that cannot be returned to the
// caller
//
// Deprecated: The whole module is deprecated. Use span.RecordError and span.SetStatus
// from go.opentelemetry.io/otel/trace and go.opentelemetry.io/otel/codes.
type ErrorHandler func(error)
