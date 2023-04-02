package jwt

import (
	"fmt"
	"notionboy/internal/pkg/config"
	"notionboy/internal/pkg/logger"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	signingMethod = jwt.SigningMethodHS256
	signingKey    string
)

func init() {
	signingKey = config.GetConfig().JWT.SigningKey
	if signingKey == "" {
		logger.SugaredLogger.Panicw("JWT SigningKey is empty")
	}
}

// GenerateToken creates a new JWT token with the given user ID as the subject.
func GenerateToken(userID string) (string, error) {
	claims := NewClaims(userID)
	token := jwt.NewWithClaims(signingMethod, claims)

	tokenString, err := token.SignedString([]byte(signingKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken parses and validates a JWT token and returns the user ID if the token is valid.
func ValidateToken(tokenString string) (string, error) {
	// Parse the token using the secret key and custom claims.
	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Check that the signing method is HMAC with SHA-256.
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Return the secret key for signature verification.
		return []byte(signingKey), nil
	})
	if err != nil {
		logger.SugaredLogger.Errorw("parse token failed", "token", tokenString, "err", err)
		return "", err
	}

	// Verify that the token is valid and has not been tampered with.
	if !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	// logger.SugaredLogger.Debugw("check token", "token", token)

	// Get the custom claims from the token.
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid token claims")
	}

	// Check the token's expiration time to ensure that it has not expired.
	exp, ok := claims["exp"].(float64)
	if !ok {
		return "", fmt.Errorf("invalid expiration time")
	}
	if time.Now().Unix() > int64(exp) {
		return "", fmt.Errorf("token has expired")
	}

	// Get the user ID from the custom claims.
	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", fmt.Errorf("invalid user ID")
	}

	return userID, nil
}

// NewClaims creates a new set of custom claims for a JWT token.
func NewClaims(userID string) jwt.Claims {
	// Get the JWT expiration time from the configuration, or use the default of 1 hour.
	expiration := config.GetConfig().JWT.Expiration
	if expiration == 0 {
		expiration = time.Hour
	}
	logger.SugaredLogger.Debugw("generate token", "userID", userID, "expiration", expiration)

	// Return the custom claims as a map.
	return jwt.MapClaims{
		"sub":     "user",
		"user_id": userID,
		"iat":     time.Now().Unix(),
		"exp":     time.Now().Add(expiration).Unix(),
	}
}
