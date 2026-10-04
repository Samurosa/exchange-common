package session

import (
	"context"
	"errors"
)

var ErrNotFound = errors.New("session not found")

type Validator interface {
	Validate(ctx context.Context, jti string) error
}
