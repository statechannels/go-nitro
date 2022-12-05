package internal

import (
	"fmt"
)

type wrappableError struct {
	message string
	cause   error
}

var _ error = (*wrappableError)(nil)

func NewError(message string) error {
	return &wrappableError{
		message: message,
		cause:   nil,
	}
}

func WrapError(err error, cause error) error {
	return &wrappableError{
		message: err.Error(),
		cause:   cause,
	}
}

func (e *wrappableError) Error() string {
	if e.cause == nil {
		return e.message
	}

	return fmt.Sprintf("%s: %s", e.message, e.cause.Error())
}

func (e *wrappableError) Unwrap() error {
	return e.cause
}

func (e *wrappableError) Is(target error) bool {
	return e.message == target.Error()
}

func (e *wrappableError) As(target interface{}) bool {
	if t, ok := target.(*wrappableError); ok {
		*t = *e

		return true
	}

	return false
}
