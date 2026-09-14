package jwt

import (
	"github.com/golang-jwt/jwt/v5"
)

type Parser struct {
	secret []byte
}

func NewParser(secret string) *Parser {
	return &Parser{
		secret: []byte(secret),
	}
}

func (j *Parser) ParseAccessToken(accessToken string) (AccessTokenParse, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&AccessTokenParse{},
		func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, ErrInvalidToken
			}
			return j.secret, nil
		},
	)
	if err != nil {
		return AccessTokenParse{}, ErrInvalidToken
	}

	claims, err := GetClaimsWithAccessToken(token)
	if err != nil {
		return AccessTokenParse{}, err
	}

	return *claims, nil
}

func GetClaimsWithAccessToken(token *jwt.Token) (*AccessTokenParse, error) {

	if !token.Valid {
		return &AccessTokenParse{}, ErrInvalidToken
	}

	claims, ok := token.Claims.(*AccessTokenParse)
	if !ok {
		return &AccessTokenParse{}, ErrInvalidToken
	}

	return &AccessTokenParse{
		UserID: claims.UserID,
		Role:   claims.Role,
		Device: claims.Device,
		JTI:    claims.JTI,
	}, nil
}
