package svc

import "errors"

var (
	ErrUserExists    = errors.New("user already exists")
	ErrEntryNotFound = errors.New("the provided entry could not be found")
)
