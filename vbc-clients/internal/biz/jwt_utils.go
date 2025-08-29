package biz

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"vbc/lib"
)

func JWTParse(tokenStr string, secret string) (lib.TypeMap, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return lib.TypeMap(claims), nil
	}
	return nil, nil
}
