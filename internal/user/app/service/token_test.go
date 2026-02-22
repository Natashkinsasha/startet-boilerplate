//go:build unit

package service

import (
	"testing"
	"time"

	"starter-boilerplate/pkg/jwt"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testTokenService() TokenService {
	mgr := jwt.NewManager(jwt.Config{
		AccessSecret:  "test-access-secret",
		RefreshSecret: "test-refresh-secret",
		AccessTTL:     15 * time.Minute,
		RefreshTTL:    24 * time.Hour,
	})
	return NewTokenService(mgr)
}

func TestIssueTokenPair_Success(t *testing.T) {
	svc := testTokenService()

	pair, err := svc.IssueTokenPair("user-1", "admin")

	require.NoError(t, err)
	assert.NotEmpty(t, pair.AccessToken)
	assert.NotEmpty(t, pair.RefreshToken)
}

func TestValidateRefreshToken_Success(t *testing.T) {
	svc := testTokenService()

	pair, err := svc.IssueTokenPair("user-1", "admin")
	require.NoError(t, err)

	claims, err := svc.ValidateRefreshToken(pair.RefreshToken)

	require.NoError(t, err)
	assert.Equal(t, "user-1", claims.UserID)
}

func TestValidateRefreshToken_Invalid(t *testing.T) {
	svc := testTokenService()

	_, err := svc.ValidateRefreshToken("garbage-token")

	assert.EqualError(t, err, "invalid refresh token")
}
