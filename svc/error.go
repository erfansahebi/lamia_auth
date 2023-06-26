package svc

import "errors"

var (
	ErrUserExists        = errors.New("user already exists")
	ErrUserDoesNotExists = errors.New("user doesn't exists")
	ErrEntryNotFound     = errors.New("the provided entry could not be found")
)
