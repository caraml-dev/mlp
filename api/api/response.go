package api

import (
	"encoding/json"
	"errors"
	"net/http"

	apperror "github.com/caraml-dev/mlp/api/pkg/errors"
)

type Response struct {
	code int
	data interface{}
}

type ErrorMessage struct {
	Message string `json:"error"`
}

func (r *Response) WriteTo(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(r.code)

	if r.data != nil {
		encoder := json.NewEncoder(w)
		_ = encoder.Encode(r.data)
	}
}

func Ok(data interface{}) *Response {
	return &Response{
		code: http.StatusOK,
		data: data,
	}
}

func Created(data interface{}) *Response {
	return &Response{
		code: http.StatusCreated,
		data: data,
	}
}

func NoContent() *Response {
	return &Response{
		code: http.StatusNoContent,
	}
}

func Error(code int, msg string) *Response {
	return &Response{
		code: code,
		data: ErrorMessage{msg},
	}
}

func NotFound(msg string) *Response {
	return Error(http.StatusNotFound, msg)
}

func BadRequest(msg string) *Response {
	return Error(http.StatusBadRequest, msg)
}

func InternalServerError(msg string) *Response {
	return Error(http.StatusInternalServerError, msg)
}

func Forbidden(msg string) *Response {
	return Error(http.StatusForbidden, msg)
}

func Conflict(msg string) *Response {
	return Error(http.StatusConflict, msg)
}

func FromError(err error) *Response {
	if errors.Is(err, &apperror.NotFoundError{}) {
		return NotFound(err.Error())
	} else if errors.Is(err, &apperror.AlreadyExistsError{}) {
		return Conflict(err.Error())
	} else if errors.Is(err, &apperror.InvalidArgumentError{}) {
		return BadRequest(err.Error())
	}

	return InternalServerError(err.Error())
}
