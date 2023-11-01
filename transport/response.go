package transport

import (
	"auth-service/endpoint"
	"auth-service/language"
	"net/http"
)

type MetaResponse struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

type Response struct {
	Meta *MetaResponse `json:"meta,omitempty"`
	Data interface{}   `json:"data,omitempty"`
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

	if response == nil {
		code = http.StatusInternalServerError
		message = localizer.Get("failed")
		return code, Response{
			Meta: &MetaResponse{
				Code:    code,
				Message: message,
			},
		}
	}

	switch valueType := response.(type) {
	case *endpoint.ErrorResponse:
		code = valueType.StatusCode
		if code == 0 {
			code = endpoint.MappingErrorToCode(valueType.Err)
		}
		message = localizer.MapError(valueType.Err)

	case error:
		code = endpoint.MappingErrorToCode(valueType)
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
