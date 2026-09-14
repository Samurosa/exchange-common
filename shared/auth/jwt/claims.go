package jwt

import "github.com/golang-jwt/jwt/v5"

type AccessClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	Device string `json:"device"`
	JTI    string `json:"jti"`

	jwt.RegisteredClaims
}

type RefreshClaims struct {
	AccessTokenJTI  string `json:"access_token_jti"`
	RefreshTokenJTI string `json:"refresh_token_jti"`

	jwt.RegisteredClaims
}
