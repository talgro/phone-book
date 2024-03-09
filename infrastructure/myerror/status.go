package myerror

import (
	"fmt"
)

type errorType int

const (
	InternalError errorType = iota
	BadRequestError
	NotFoundError
	ForbiddenError
)

type MyError struct {
	Message string
	Type    errorType
}

func (e MyError) Error() string {
	return e.Message
}

func Wrap(err error, format string, args ...interface{}) error {
	parsedErr := GetParsedError(err)
	message := fmt.Sprintf(format, args...)

	return MyError{
		Message: fmt.Sprintf("%s: %s", message, parsedErr.Error()),
		Type:    parsedErr.Type,
	}
}

func NewInternalError(format string, args ...interface{}) error {
	message := fmt.Sprintf(format, args...)
	return MyError{
		Message: message,
		Type:    InternalError,
	}
}

func NewBadRequestError(format string, args ...interface{}) error {
	message := fmt.Sprintf(format, args...)
	return MyError{
		Message: message,
		Type:    BadRequestError,
	}
}

func NewNotFoundError(format string, args ...interface{}) error {
	message := fmt.Sprintf(format, args...)
	return MyError{
		Message: message,
		Type:    NotFoundError,
	}
}

func NewForbiddenError(format string, args ...interface{}) error {
	message := fmt.Sprintf(format, args...)
	return MyError{
		Message: message,
		Type:    ForbiddenError,
	}
}

func GetParsedError(err error) MyError {
	if e, ok := err.(MyError); ok {
		return e
	}

	message := err.Error()
	if err == nil {
		message = "unknown error"
	}

	return MyError{
		Message: message,
		Type:    InternalError,
	}
}
