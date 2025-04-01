package util

import "errors"

var (
	ErrNotFound  = errors.New("Resource not found.")
	ErrConflict  = errors.New("Resource conflict.")
	ErrMalformed = errors.New("Bad request.")
	ErrForbidden = errors.New("Forbidden to preform action.")
	ErrInternal  = errors.New("Unexpected internal error.")
)
