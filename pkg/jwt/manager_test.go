//go:build unit

package jwt

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testManager() *Manager {
	return NewManager(Config{
		AccessSecret:  "test-access-secret",
		RefreshSecret: "test-refresh-secret",
		AccessTTL:     15 * time.Minute,
		RefreshTTL:    24 * time.Hour,
	})
}

func TestGenerateAccessToken(t *testing.T) {
	m := testManager()

	token, err := m.GenerateAccessToken("user-1", "admin")

	require.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestGenerateRefreshToken(t *testing.T) {
	m := testManager()

	token, err := m.GenerateRefreshToken("user-1", "admin")

	require.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestValidateAccessToken(t *testing.T) {
	m := testManager()

	token, err := m.GenerateAccessToken("user-1", "admin")
	require.NoError(t, err)

	claims, err := m.ValidateAccessToken(token)

	require.NoError(t, err)
	assert.Equal(t, "user-1", claims.UserID)
	assert.Equal(t, "admin", claims.Role)
}

func TestValidateRefreshToken(t *testing.T) {
	m := testManager()

	token, err := m.GenerateRefreshToken("user-1", "user")
	require.NoError(t, err)

	claims, err := m.ValidateRefreshToken(token)

	require.NoError(t, err)
	assert.Equal(t, "user-1", claims.UserID)
	assert.Equal(t, "user", claims.Role)
}

func TestAccessTokenCannotBeValidatedAsRefresh(t *testing.T) {
	m := testManager()

	token, err := m.GenerateAccessToken("user-1", "admin")
	require.NoError(t, err)

	_, err = m.ValidateRefreshToken(token)

	assert.Error(t, err)
}

func TestRefreshTokenCannotBeValidatedAsAccess(t *testing.T) {
	m := testManager()

	token, err := m.GenerateRefreshToken("user-1", "admin")
	require.NoError(t, err)

	_, err = m.ValidateAccessToken(token)

	assert.Error(t, err)
}

func TestInvalidTokenString(t *testing.T) {
	m := testManager()

	_, err := m.ValidateAccessToken("garbage.token.string")
	assert.Error(t, err)

	_, err = m.ValidateRefreshToken("garbage.token.string")
	assert.Error(t, err)
}
