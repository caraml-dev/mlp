package api

import (
	"encoding/json"
	"net/http"
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
