package jwt

import "github.com/golang-jwt/jwt/v5"

type AccessTokenClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	Device string `json:"device"`

	jwt.RegisteredClaims
}

func (c AccessTokenClaims) Validate() error {
	if c.UserID == "" {
		return ErrInvalidToken
	}

	if c.Role == "" {
		return ErrInvalidToken
	}

	if c.Device == "" {
		return ErrInvalidToken
	}

	if c.ID == "" {
		return ErrInvalidToken
	}

	if c.ExpiresAt == nil {
		return ErrInvalidToken
	}

	if c.IssuedAt == nil {
		return ErrInvalidToken
	}

	return nil
}
