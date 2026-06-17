package auth

import (
	"errors"
	"fmt"

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
		// Verify the signing method to prevent algorithm confusion attacks
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
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
