package errors

import "context"

// ErrorHandler allows for handling of error that cannot be returned to the
// caller
type ErrorHandler func(context.Context, error) error
