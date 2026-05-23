package auth

import (
	"errors"

	"github.com/golang-jwt/jwt"
	"github.com/tapiaw38/practiq-be/internal/platform/config"
)

type CustomClaims struct {
	UserID       string      `json:"user_id"`
	TokenVersion uint        `json:"token_version"`
	Roles        []RoleClaim `json:"roles"`
	jwt.StandardClaims
}

type RoleClaim struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func ValidateToken(tokenStr string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		secret := config.GetConfigService().ServerConfig.JWTSecret
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}
