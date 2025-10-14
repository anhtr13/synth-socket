package util

import (
	"fmt"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"

	"github.com/anhtr13/synth-socket/api/internal/conf"
)

type JwtClaims struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	jwt.RegisteredClaims
}

func VerifyJWT(tokenString string) (map[string]any, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(conf.JWT_SEC), nil
	})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	claims := token.Claims.(jwt.MapClaims)
	exp, err := claims.GetExpirationTime()
	if err != nil {
		return nil, fmt.Errorf("invalid token")
	}
	if exp.Compare(time.Now()) < 0 {
		return nil, fmt.Errorf("token expired")
	}
	return claims, nil
}

func SignJWT(payload JwtClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	tokenString, err := token.SignedString([]byte(conf.JWT_SEC))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
