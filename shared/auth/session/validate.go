package session

import "context"

type Validator interface {
	Validate(ctx context.Context, jti string) error
}
