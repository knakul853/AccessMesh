package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJWTTokenFlow(t *testing.T) {
	// Test token generation
	role := "admin"
	token, err := GenerateToken(role)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Test token validation
	claims, err := ValidateToken(token)
	assert.NoError(t, err)
	assert.Equal(t, role, claims.Role)
	assert.True(t, claims.ExpiresAt.Time.After(time.Now()))
}

func TestInvalidToken(t *testing.T) {

	_, err := ValidateToken("")
	assert.Error(t, err)

	_, err = ValidateToken("invalid.token.string")
	assert.Error(t, err)

	expiredToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlIjoiYWRtaW4iLCJleHAiOjE1MTYyMzkwMjJ9.1234567890"
	_, err = ValidateToken(expiredToken)
	assert.Error(t, err)
}

// func setupTestStore(t *testing.T) *store.MongoStore {
// 	store, err := store.NewMongoStore("mongodb://localhost:27017/pbac_test")
// 	if err != nil {
// 		t.Fatalf("Failed to connect to test database: %v", err)
// 	}
// 	return store
// }
