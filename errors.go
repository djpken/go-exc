package exc

import (
	"errors"

	"github.com/djpken/go-exc/types"
)

var (
	// ErrInvalidExchange is returned when an invalid exchange type is specified
	ErrInvalidExchange = errors.New("invalid exchange type")

	// ErrInvalidConfig is returned when the configuration is invalid
	ErrInvalidConfig = errors.New("invalid configuration")

	// ErrNotConnected is returned when attempting an operation while not connected
	ErrNotConnected = errors.New("not connected")

	// ErrAlreadyConnected is returned when attempting to connect while already connected
	ErrAlreadyConnected = errors.New("already connected")

	// ErrInvalidSymbol is returned when an invalid symbol is specified
	ErrInvalidSymbol = errors.New("invalid symbol")

	// ErrInvalidOrder is returned when an invalid order is specified
	ErrInvalidOrder = errors.New("invalid order")

	// ErrOrderNotFound is returned when an order is not found
	ErrOrderNotFound = errors.New("order not found")

	// ErrInsufficientBalance is returned when there is insufficient balance
	ErrInsufficientBalance = errors.New("insufficient balance")

	// ErrRateLimitExceeded is returned when rate limit is exceeded
	ErrRateLimitExceeded = errors.New("rate limit exceeded")

	// ErrNotImplemented is returned when a feature is not implemented
	ErrNotImplemented = errors.New("not implemented")

	// ErrNotSupported is returned when a feature is not supported by the exchange
	// This is an alias to types.ErrNotSupported for convenience
	ErrNotSupported = types.ErrNotSupported
)

// Error represents an exchange API error
type Error struct {
	// Code is the error code from the exchange
	Code string

	// Message is the error message
	Message string

	// Err is the underlying error
	Err error
}

// Error implements the error interface
func (e *Error) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

// Unwrap returns the underlying error
func (e *Error) Unwrap() error {
	return e.Err
}

// NewError creates a new Error
func NewError(code, message string, err error) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Err:     err,
	}
}
