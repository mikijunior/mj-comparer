package store

import "errors"

var (
	ErrorRecordNotFound    = errors.New("Record not found")
	ErrUnauthorized        = errors.New("Unauthorized")
	ErrInternalServerError = errors.New("Something went wrong")
)
