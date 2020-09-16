package api

import (
	"encoding/json"
	"net/http"
)

type ApiResponse struct {
	code int
	data interface{}
}

type ErrorMessage struct {
	Message string `json:"error"`
}

func (r *ApiResponse) WriteTo(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(r.code)

	if r.data != nil {
		encoder := json.NewEncoder(w)
		encoder.Encode(r.data)
	}
}

func Ok(data interface{}) *ApiResponse {
	return &ApiResponse{
		code: http.StatusOK,
		data: data,
	}
}

func Created(data interface{}) *ApiResponse {
	return &ApiResponse{
		code: http.StatusCreated,
		data: data,
	}
}

func NoContent() *ApiResponse {
	return &ApiResponse{
		code: http.StatusNoContent,
	}
}

func Error(code int, msg string) *ApiResponse {
	return &ApiResponse{
		code: code,
		data: ErrorMessage{msg},
	}
}

func NotFound(msg string) *ApiResponse {
	return Error(http.StatusNotFound, msg)
}

func BadRequest(msg string) *ApiResponse {
	return Error(http.StatusBadRequest, msg)
}

func InternalServerError(msg string) *ApiResponse {
	return Error(http.StatusInternalServerError, msg)
}

func Forbidden(msg string) *ApiResponse {
	return Error(http.StatusForbidden, msg)
}
