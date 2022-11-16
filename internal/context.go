package internal

import "context"

// TraceContextKey for tracer context metadata
type TraceContextKey struct{}

// ExtendedContextWithMetadata add metadata to base context
func ExtendedContextWithMetadata[metaType any](baseCtx context.Context, metaKey any, metadata metaType) context.Context {
	return context.WithValue(baseCtx, metaKey, metadata)
}

// GetContextMetadata will try get form context.Context metadata
func GetContextMetadata[metaType any](baseCtx context.Context, metaKey any) (metaData metaType, isExist bool) {
	metaData, isExist = baseCtx.Value(metaKey).(metaType)

	return
}
