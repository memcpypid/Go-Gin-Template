package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenClaims struct {
	UserID uuid.UUID `json:"user_id"`
	Role   string    `json:"role"`
	jwt.RegisteredClaims
}

func GenerateTokens(userID uuid.UUID, role, secret string, accessExpHours, refreshExpHours int) (accessToken, refreshToken string, refreshExpAt time.Time, err error) {
	secretKey := []byte(secret)

	// Access Token
	accessClaims := TokenClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(accessExpHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	accessTokenRaw := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenStr, err := accessTokenRaw.SignedString(secretKey)
	if err != nil {
		return "", "", time.Time{}, err
	}

	// Refresh Token
	refreshExpAt = time.Now().Add(time.Duration(refreshExpHours) * time.Hour)
	refreshClaims := jwt.RegisteredClaims{
		Subject:   userID.String(),
		ExpiresAt: jwt.NewNumericDate(refreshExpAt),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	refreshTokenRaw := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenStr, err := refreshTokenRaw.SignedString(secretKey)
	if err != nil {
		return "", "", time.Time{}, err
	}

	return accessTokenStr, refreshTokenStr, refreshExpAt, nil
}

func ValidateToken(tokenString, secret string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
