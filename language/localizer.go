package language

import (
	"database/sql"
	"errors"
)

type Localizer map[string]string

func (s Localizer) Get(key string) string {
	value, ok := s[key]
	if !ok {
		return ""
	}
	return value
}

func (s Localizer) GetError(key string) error {
	value := s.Get(key)
	return errors.New(value)
}

func (s Localizer) MapError(err error) string {
	// Get error code for each error
	var errCode string
	{
		switch err {
		case sql.ErrNoRows:
			errCode = "not_found"
		default:
			_, ok := s[err.Error()]
			if ok {
				errCode = err.Error()
			}
		}
	}

	// Get error message
	value, ok := s[errCode]
	if ok {
		return value
	} else {
		return err.Error()
	}
}
