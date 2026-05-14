package crypt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaim holds the data embedded inside a JWT token.
type JWTClaim struct {
	UserID   uint     `json:"user_id"`
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Roles    []string `json:"roles"`
	jwt.RegisteredClaims
}

// GenerateJWT creates a signed JWT token for the given user.
func GenerateJWT(userID uint, username string, email string, roles []string, secret string, expDays int) (string, error) {
	expirationTime := time.Now().Add(time.Duration(expDays) * 24 * time.Hour)
	claims := &JWTClaim{
		UserID:   userID,
		Username: username,
		Email:    email,
		Roles:    roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	return tokenString, err
}

// ParseJWT validates a token string and returns the embedded claims.
func ParseJWT(tokenString string, secret string) (*JWTClaim, error) {
	claims := &JWTClaim{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
