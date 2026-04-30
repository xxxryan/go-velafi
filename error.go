package velafi

import (
	"errors"
	"fmt"
)

type Error struct {
	HTTPStatus int
	Code       int
	Message    string
	RequestID  string
}

func (e *Error) Error() string {
	if e.RequestID != "" {
		return fmt.Sprintf("velafi: %d %s (http %d, request_id=%s)", e.Code, e.Message, e.HTTPStatus, e.RequestID)
	}
	return fmt.Sprintf("velafi: %d %s (http %d)", e.Code, e.Message, e.HTTPStatus)
}

func IsNotFound(err error) bool {
	var e *Error
	return errors.As(err, &e) && e.HTTPStatus == 404
}

func IsUnauthorized(err error) bool {
	var e *Error
	return errors.As(err, &e) && e.HTTPStatus == 401
}

func IsRateLimited(err error) bool {
	var e *Error
	return errors.As(err, &e) && e.HTTPStatus == 429
}
