package custom

import (
	"errors"
)

var (
	ErrAlreadyExists = errors.New("resource already exists")
	ErrNotFound      = errors.New("resource not found")
	ErrInternal      = errors.New("internal server error")
	ErrUnauthorized  = errors.New("unauthorized")
	ErrForbidden     = errors.New("you are not authorized to access this resource")
)
