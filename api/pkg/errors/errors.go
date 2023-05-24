package errors

import "fmt"

// NotFoundError is an error type that indicates that the resource is not found
type NotFoundError struct {
	message string
}

// NewNotFoundErrorf creates a new NotFoundError with a formatted message
func NewNotFoundErrorf(format string, a ...any) *NotFoundError {
	return &NotFoundError{
		message: fmt.Sprintf(format, a...),
	}
}

// Error returns the error message
func (e *NotFoundError) Error() string {
	return e.message
}

// Is check whether the error is NotFoundError
func (e *NotFoundError) Is(target error) bool {
	_, ok := target.(*NotFoundError)
	return ok
}

// AlreadyExistsError is an error type that indicates that the resource already exists
type AlreadyExistsError struct {
	message string
}

// NewAlreadyExistsErrorf creates a new AlreadyExistsError with a formatted message
func NewAlreadyExistsErrorf(format string, a ...any) *AlreadyExistsError {
	return &AlreadyExistsError{
		message: fmt.Sprintf(format, a...),
	}
}

// Error returns the error message
func (e *AlreadyExistsError) Error() string {
	return e.message
}

// Is check whether the error is AlreadyExistsError
func (e *AlreadyExistsError) Is(target error) bool {
	_, ok := target.(*AlreadyExistsError)
	return ok
}

// InvalidArgumentError is an error type that indicates that the argument is invalid
type InvalidArgumentError struct {
	message string
}

// NewInvalidArgumentErrorf creates a new InvalidArgumentError with a formatted message
func NewInvalidArgumentErrorf(format string, a ...any) *InvalidArgumentError {
	return &InvalidArgumentError{
		message: fmt.Sprintf(format, a...),
	}
}

// Error returns the error message
func (e *InvalidArgumentError) Error() string {
	return e.message
}

// Is check whether the error is InvalidArgumentError
func (e *InvalidArgumentError) Is(target error) bool {
	_, ok := target.(*InvalidArgumentError)
	return ok
}
