package errors

import "errors"

var (
	ErrLoginTaken      = errors.New("login already taken")
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidPassword = errors.New("invalid password")
)
