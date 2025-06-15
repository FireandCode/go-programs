package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWTClaim struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateJWT(userID uint, secret string, expirationTime time.Duration) (string, error) {
	claims := &JWTClaim{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expirationTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(tokenString string, secret string) (*JWTClaim, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		return nil, errors.New("couldn't parse claims")
	}

	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("token expired")
	}

	return claims, nil
}

func GenerateAccessToken(userID uint, secret string) (string, error) {
	return GenerateJWT(userID, secret, 10*time.Minute)
}

func GenerateRefreshToken(userID uint, secret string) (string, error) {
	return GenerateJWT(userID, secret, 24*7*time.Hour) // 1 week
} 