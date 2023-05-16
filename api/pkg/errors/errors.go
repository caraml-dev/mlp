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

// Is check whether the error is NotFoundError
func (e *NotFoundError) Is(target error) bool {
	_, ok := target.(*NotFoundError)
	return ok
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

func (e *AlreadyExistsError) Is(target error) bool {
	_, ok := target.(*AlreadyExistsError)
	return ok
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

func (e *InvalidArgumentError) Is(target error) bool {
	_, ok := target.(*InvalidArgumentError)
	return ok
}
