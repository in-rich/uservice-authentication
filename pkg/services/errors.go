package services

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")

	ErrUnauthenticated  = errors.New("unauthenticated")
	ErrVerifyToken      = errors.New("verify token")
	ErrEmailNotVerified = errors.New("email not verified")

	ErrInvalidUpdateUser = errors.New("invalid update user")
)
