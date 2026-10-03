package jwt

import (
	"github.com/golang-jwt/jwt/v5"
)

type Parser struct {
	secret []byte
}

func NewParser(secret string) (*Parser, error) {
	if secret == "" {
		return nil, ErrEmptySecret
	}

	return &Parser{
		secret: []byte(secret),
	}, nil
}

func (j *Parser) ParseAccessToken(accessToken string) (AccessTokenClaims, error) {
	if accessToken == "" {
		return AccessTokenClaims{}, ErrInvalidToken
	}

	token, err := jwt.ParseWithClaims(
		accessToken,
		&AccessTokenClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, ErrInvalidToken
			}

			return j.secret, nil
		},
	)
	if err != nil {
		return AccessTokenClaims{}, err
	}

	claims, err := GetClaimsWithAccessToken(token)
	if err != nil {
		return AccessTokenClaims{}, err
	}

	return *claims, nil
}

func GetClaimsWithAccessToken(
	token *jwt.Token,
) (*AccessTokenClaims, error) {
	if token == nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*AccessTokenClaims)
	if !ok || claims == nil {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
