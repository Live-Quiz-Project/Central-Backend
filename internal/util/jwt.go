package util

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type Claims struct {
	UserID       uuid.UUID `json:"uid"`
	Name         string    `json:"name"`
	DisplayName  string    `json:"display_name"`
	DisplayEmoji string    `json:"display_emoji"`
	DisplayColor string    `json:"display_color"`
	jwt.StandardClaims
}

func GenerateToken(uid uuid.UUID, duration time.Time, secret string) (string, error) {
	claims := &Claims{
		UserID: uid,
		StandardClaims: jwt.StandardClaims{
			Issuer:    uid.String(),
			ExpiresAt: duration.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func DecodeToken(tokenString string, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, err
	}

	return claims, nil
}
