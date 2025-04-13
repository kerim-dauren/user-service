package domain

import "errors"

var (
	// ErrUserNotFound will throw if the requested user is not exists
	ErrUserNotFound          = errors.New("user not found")
	ErrUserMailAlreadyExists = errors.New("user mail already exists")
)
