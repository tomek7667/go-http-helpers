package utils

import (
	"github.com/golang-jwt/jwt"
)

func JwtEncode(claims map[string]any, secret string) (string, error) {
	t := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims(claims),
	)
	return t.SignedString([]byte(secret))
}

func JwtVerify(token, secret string) (map[string]any, error) {
	t, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	return t.Claims.(jwt.MapClaims), nil
}
