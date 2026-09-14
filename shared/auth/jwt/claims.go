package jwt

import "github.com/golang-jwt/jwt/v5"

type AccessTokenParse struct {
	UserID string
	Role   string
	Device string
	JTI    string

	jwt.RegisteredClaims
}
