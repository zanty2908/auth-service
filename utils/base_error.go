package utils

import "errors"

var (
	ErrorFailed       = errors.New("error")
	ErrorNotFound     = errors.New("not_found")
	ErrorExist        = errors.New("exist")
	ErrorInvalidPhone = errors.New("invalid_phone")
	ErrorInvalidMail  = errors.New("invalid_email")
	ErrorBadRequest   = errors.New("bad_request")
	ErrorMissingId    = errors.New("missing_id")
	ErrorBlocked      = errors.New("user_blocked")
)
