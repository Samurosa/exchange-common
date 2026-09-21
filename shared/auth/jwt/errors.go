package jwt

import "errors"

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrEmptySecret  = errors.New("jwt secret is empty")
)
