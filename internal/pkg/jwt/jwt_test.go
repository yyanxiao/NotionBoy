package jwt

import (
	"notionboy/internal/pkg/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

// test GenerateToken

func TestJWT(t *testing.T) {
	// Set the JWT signing key in the config
	config.GetConfig().JWT.SigningKey = "secret"

	// Generate a token
	userID := "123"
	token, err := GenerateToken(userID)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Validate the token
	validUserID, err := ValidateToken(token)
	assert.NoError(t, err)
	assert.Equal(t, userID, validUserID)
}
