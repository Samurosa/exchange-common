package sharedContext

import "errors"

var (
	ErrInvalidContext       = errors.New("invalid context")
	ErrCtxForUpdateNotFound = errors.New("context for update not found")
)
