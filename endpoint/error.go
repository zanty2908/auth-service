package endpoint

import (
	"errors"
	"fmt"
	"net/http"
)

var errorMap = map[string]int{
	"invalid_phone":           http.StatusBadRequest,
	"invalid_email":           http.StatusBadRequest,
	"missing_id":              http.StatusBadRequest,
	"bad_request":             http.StatusBadRequest,
	"unauthorized_client":     http.StatusUnauthorized,
	"access_denied":           http.StatusForbidden,
	"server_error":            http.StatusInternalServerError,
	"temporarily_unavailable": http.StatusServiceUnavailable,
}

var errorGRPCMap = map[string]int{
	"invalid_phone":           http.StatusBadRequest,
	"invalid_email":           http.StatusBadRequest,
	"missing_id":              http.StatusBadRequest,
	"bad_request":             http.StatusBadRequest,
	"unauthorized_client":     http.StatusUnauthorized,
	"access_denied":           http.StatusForbidden,
	"server_error":            http.StatusInternalServerError,
	"temporarily_unavailable": http.StatusServiceUnavailable,
}

type ErrorResponse struct {
	StatusCode int
	Err        error
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf(`{"code":%d,"message":"%v"}`, r.StatusCode, r.Err)
}

func MappingErrorToCode(err error) int {
	code, ok := errorMap[err.Error()]
	if !ok {
		return http.StatusInternalServerError
	}
	return code
}

func MappingErrorGRPCToCode(err error) int {
	code, ok := errorGRPCMap[err.Error()]
	if !ok {
		return http.StatusInternalServerError
	}
	return code
}

func NewErrorRes(code int, err error) *ErrorResponse {
	return &ErrorResponse{StatusCode: code, Err: err}
}

func BadRequestError() *ErrorResponse {
	return NewErrorRes(http.StatusBadRequest, errors.New("bad_request"))
}

func InternalServerError() *ErrorResponse {
	return NewErrorRes(http.StatusInternalServerError, errors.New("server_error"))
}
