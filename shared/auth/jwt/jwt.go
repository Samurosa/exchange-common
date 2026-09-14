package jwt

import "github.com/golang-jwt/jwt/v5"

type Parser struct {
	secret []byte
}

func NewParser(secret string) *Parser {
	return &Parser{
		secret: []byte(secret),
	}
}

func (p *Parser) ParseAccessToken(tokenString string) (AccessClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&AccessClaims{},
		func(token *jwt.Token) (any, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, ErrInvalidToken
			}

			return p.secret, nil
		},
	)
	if err != nil {
		return AccessClaims{}, ErrInvalidToken
	}

	claims, ok := token.Claims.(*AccessClaims)
	if !ok || !token.Valid {
		return AccessClaims{}, ErrInvalidToken
	}

	return *claims, nil
}
