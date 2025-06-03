package validator

import "errors"

var (
	ErrUserOrEmailAlreadyExists = errors.New("user or email already exists")
	ErrInvalidEmail             = errors.New("must be a valid email address")
	ErrNotBlank                 = errors.New("must be provided")
	ErrInvalidCredentials       = errors.New("invalid credentials provided")
)
