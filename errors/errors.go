package errors

// ErrorHandler allows for handling of error that cannot be returned to the
// caller
type ErrorHandler func(error)
