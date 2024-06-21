package internal

import "fmt"

// Error is a custom error type that supports wrapping and error codes
// It includes the original error, a message and an error code
type Error struct {
	code     ErrorCode
	original error
	message  string
}

// ErrorCode defines the supported error codes
type ErrorCode uint

const (
	ErrUnknown ErrorCode = iota
	ErrNotFound
	ErrUniqueConstraint
	ErrInvalidInput
)

// WrapErrorf returns a new error that wraps the original error and includes a message and an error code
func WrapErrorf(original error, code ErrorCode, format string, args ...interface{}) error {
	return &Error{
		code:     code,
		original: original,
		message:  fmt.Sprintf(format, args...),
	}
}

// NewErrorf instantiates a new error
func NewErrorf(code ErrorCode, format string, args ...interface{}) error {
	return WrapErrorf(nil, code, format, args...)
}

// Error returns the error message that includes the original error if present
func (e *Error) Error() string {
	if e.original != nil {
		return fmt.Sprintf("%s: %v", e.message, e.original)
	}

	return e.message
}

// Unwrap returns the original error
func (e *Error) Unwrap() error {
	return e.original
}

// Code returns the error code
func (e *Error) Code() ErrorCode {
	return e.code
}
