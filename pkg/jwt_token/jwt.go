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

var MapTypeToken = map[string]time.Duration{
	"token":         time.Hour * 3,
	"refresh_token": time.Hour * 72,
}

var jwtSecret = []byte(env.GetEnv("APP_SECRET", ""))

func GenerateToken(ctx context.Context, username, fullname string, tokenType string, now time.Time) (string, error) {
	claimToken := ClaimToken{
		Username: username,
		FullName: fullname,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    env.GetEnv("APP_NAME", ""),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(MapTypeToken[tokenType])),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimToken)
	resultToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return resultToken, fmt.Errorf("failed to generate claims token: %v", err)
	}

	return resultToken, nil
}

func ValidateToken(ctx context.Context, token string) (*ClaimToken, error) {
	var (
		claimToken *ClaimToken
		ok         bool
	)

	jwtToken, err := jwt.ParseWithClaims(token, &ClaimToken{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("failed to validate token: %v", t.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse jwt token: %v", err)
	}

	if claimToken, ok = jwtToken.Claims.(*ClaimToken); !ok || !jwtToken.Valid {
		return nil, fmt.Errorf("token invalid")
	}

	return claimToken, nil
}
