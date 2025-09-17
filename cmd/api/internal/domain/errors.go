package domain

import "errors"

var (
	ErrToolNotFound = errors.New("tool not found")
	ErrValidation   = errors.New("validation failed")
	ErrUserNotFound = errors.New("user not found")
)
