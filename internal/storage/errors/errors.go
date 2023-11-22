package errors

import "errors"

var (
	ErrLoginTaken            = errors.New("login already taken")
	ErrUserNotFound          = errors.New("user not found")
	ErrInvalidPassword       = errors.New("invalid password")
	ErrUploadedByThisUser    = errors.New("order has already been uploaded by this user")
	ErrUploadedByAnotherUser = errors.New("order has already been uploaded by another user")
)
