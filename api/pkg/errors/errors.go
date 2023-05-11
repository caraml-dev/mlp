package errors

import "fmt"

type NotFoundError struct {
	message string
}

func NewNotFoundErrorf(format string, a ...any) *NotFoundError {
	return &NotFoundError{
		message: fmt.Sprintf(format, a...),
	}
}

func (e *NotFoundError) Error() string {
	return e.message
}

type AlreadyExistsError struct {
	message string
}

func NewAlreadyExistsErrorf(format string, a ...any) *AlreadyExistsError {
	return &AlreadyExistsError{
		message: fmt.Sprintf(format, a...),
	}
}

func (e *AlreadyExistsError) Error() string {
	return e.message
}

type InvalidArgumentError struct {
	message string
}

func NewInvalidArgumentErrorf(format string, a ...any) *InvalidArgumentError {
	return &InvalidArgumentError{
		message: fmt.Sprintf(format, a...),
	}
}

func (e *InvalidArgumentError) Error() string {
	return e.message
}
