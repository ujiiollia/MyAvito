package dbErrors

import "errors"

var (
	ErrorTagNotFound     = errors.New("tag not found")
	ErrorTagIsExist      = errors.New("tag is exist")
	ErrorFeatureNotFound = errors.New("feature not found")
	ErrorFeatureIsExist  = errors.New("feature is exist")
	ErrorUserNotFound    = errors.New("user not found")
	ErrorUserIsExist     = errors.New("user is exist")
	ErrorUserIsNotAdmin  = errors.New("user is not")
)
