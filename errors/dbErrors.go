package dberrors

import "errors"

var (
	ErrUserNotFound = errors.New("user with tis userId not found")
)
