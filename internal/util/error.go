package util

import "errors"

var (
	ErrNotFound  = errors.New("Resource not found.")
	ErrConflict  = errors.New("Resource conflict.")
	ErrMalformed = errors.New("Bad request.")
	ErrInternal  = errors.New("Unexpected internal error.")
)
