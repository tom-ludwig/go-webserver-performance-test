package utils

import (
	"go-webserver-performance-test/models"
	"time"

	"github.com/golang-jwt/jwt"
)

var secretKey = []byte("your-secret-key") // TODO: Get from env

// CreateTokens creates a pair of JWT tokens for a given username(access, refresh, error)
func CreateTokens(userID string) (string, string, error) {
	// Access token
	expirationTime := time.Now().Add(15 * time.Minute)
	accessClaims := &models.Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Issuer:    "go-webserver-performance-test",
		},
	}

	// HACK: Should use  asemetric
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	access, err := accessToken.SignedString(secretKey)
	if err != nil {
		return "", "", err
	}

	// Refresh token
	refreshExpirationTime := time.Now().Add(24 * time.Hour)
	refreshClaims := &models.Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: refreshExpirationTime.Unix(),
			Issuer:    "go-webserver-performance-test",
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(secretKey)
	if err != nil {
		return "", "", err
	}

	return access, refreshTokenString, nil
}

func DecodeToken(tokenString string) (*models.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(*models.Claims)
	if !ok {
		return nil, err
	}

	return claims, nil
}
