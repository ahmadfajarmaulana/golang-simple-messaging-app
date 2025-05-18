package jwt_token

import (
	"context"
	"fmt"
	"time"

	"simple-messaging-app/pkg/env"

	"github.com/golang-jwt/jwt/v5"
)

type ClaimToken struct {
	Username string `json:"username"`
	FullName string `json:"full_name"`
	jwt.RegisteredClaims
}

var mapTypeToken = map[string]time.Duration{
	"token":         time.Hour * 3,
	"refresh_token": time.Hour * 72,
}

func GenerateToken(ctx context.Context, username, fullname string, tokenType string) (string, error) {
	secret := env.GetEnv("APP_SECRET", "")
	claimToken := ClaimToken{
		Username: username,
		FullName: fullname,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    env.GetEnv("APP_NAME", ""),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(mapTypeToken[tokenType])),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimToken)

	resultToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return resultToken, fmt.Errorf("failed to generate claims token: %v", err)
	}

	return resultToken, nil
}
