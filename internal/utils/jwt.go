package utils

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenClaims struct {
	UserID uuid.UUID `json:"user_id"`
	Role   string    `json:"role"`
	jwt.RegisteredClaims
}

var durationRegex = regexp.MustCompile(`^(\d+)([hmd])$`)

// ParseDuration converts a duration string like "1h", "30m", or "7d" into time.Duration.
func ParseDuration(s string) (time.Duration, error) {
	match := durationRegex.FindStringSubmatch(s)
	if len(match) != 3 {
		// Fallback to standard ParseDuration if no unit or unknown unit
		return time.ParseDuration(s)
	}

	value, err := strconv.Atoi(match[1])
	if err != nil {
		return 0, err
	}

	unit := match[2]
	switch unit {
	case "m":
		return time.Duration(value) * time.Minute, nil
	case "h":
		return time.Duration(value) * time.Hour, nil
	case "d":
		return time.Duration(value) * 24 * time.Hour, nil
	default:
		return 0, fmt.Errorf("unknown unit: %s", unit)
	}
}

func GenerateTokens(userID uuid.UUID, role, secret string, accessDuration, refreshDuration time.Duration) (accessToken, refreshToken string, refreshExpAt time.Time, err error) {
	secretKey := []byte(secret)

	// Access Token
	accessClaims := TokenClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	accessTokenRaw := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenStr, err := accessTokenRaw.SignedString(secretKey)
	if err != nil {
		return "", "", time.Time{}, err
	}

	// Refresh Token
	refreshExpAt = time.Now().Add(refreshDuration)
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
