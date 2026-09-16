package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/tapiaw38/practiq-be/internal/platform/config"
)

type CustomClaims struct {
	UserID       string      `json:"user_id"`
	TokenVersion uint        `json:"token_version"`
	Roles        []RoleClaim `json:"roles"`
	// ReadOnly and ImpersonatorID are issued only by Practiq's superadmin
	// endpoint. They travel in the signed token, so every service can enforce
	// the restriction without trusting a browser flag.
	ReadOnly       bool   `json:"read_only,omitempty"`
	ImpersonatorID string `json:"impersonator_id,omitempty"`
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

// GenerateReadOnlyImpersonationToken creates a short-lived session for a
// platform operator to inspect a user's exact view. It intentionally carries
// no roles: the target's Practiq profile and school memberships define scope.
func GenerateReadOnlyImpersonationToken(userID string, tokenVersion uint, impersonatorID string) (string, error) {
	claims := CustomClaims{
		UserID:         userID,
		TokenVersion:   tokenVersion,
		Roles:          []RoleClaim{},
		ReadOnly:       true,
		ImpersonatorID: impersonatorID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(15 * time.Minute).Unix(),
		},
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(
		[]byte(config.GetConfigService().ServerConfig.JWTSecret),
	)
}
