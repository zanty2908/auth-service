package endpoint

import (
	"auth-service/language"
	"auth-service/utils"
	"fmt"
	"net/http"
)

type ErrorResponse struct {
	StatusCode int
	Err        error
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("status %d: err %v", r.StatusCode, r.Err)
}

func newErrorRes(code int, err error) *ErrorResponse {
	return &ErrorResponse{StatusCode: code, Err: err}
}

type MetaResponse struct {
	Code    int    `json:"Code"`
	Message string `json:"message"`
}

type Response struct {
	Meta *MetaResponse `json:"meta"`
	Data interface{}   `json:"data"`
}

type PagingResponse struct {
	Response
	TotalPage   int `json:"totalPage"`
	TotalResult int `json:"totalResult"`
}

func newSuccessRes(data interface{}) interface{} {
	if res, ok := data.(PagingResponse); ok {
		return res
	}

	return Response{Data: data}
}

func mappingResponse(localizer *language.Localizer, response interface{}) (int, interface{}) {
	var code int
	var message string

	switch valueType := response.(type) {
	case *ErrorResponse:
		code = valueType.StatusCode
		if code == 0 {
			code = mappingErrorToCode(valueType.Err)
		}
		message = localizer.MapError(valueType.Err)

	case error:
		code = mappingErrorToCode(valueType)
		message = localizer.MapError(valueType)

	case PagingResponse:
		return http.StatusOK, response

	default:
		return http.StatusOK, Response{Data: response}
	}

	return code, Response{
		Meta: &MetaResponse{
			Code:    code,
			Message: message,
		},
	}
}

func mappingErrorToCode(err error) (code int) {
	switch err {
	case utils.ErrorBadRequest,
		utils.ErrorMissingId,
		utils.ErrorInvalidMail:
		code = http.StatusBadRequest
	case utils.ErrorNotFound:
		code = http.StatusNotFound
	case utils.ErrorBlocked:
		code = http.StatusLocked
	default:
		code = http.StatusInternalServerError
	}
	return
}
