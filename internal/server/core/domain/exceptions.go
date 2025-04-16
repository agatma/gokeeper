package domain

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAuthentication = errors.New("user unauthorized")
	ErrUserConflict       = errors.New("user already exists")

	ErrPrivateDataBadFormat = errors.New("private data bad format")
	ErrPrivateDataNotFound  = errors.New("private data not found")
	ErrPrivateDataConflict  = errors.New("private data conflict")
)
