//go:build unit

package dto

import (
	"testing"

	"starter-boilerplate/internal/user/domain/model"

	"github.com/stretchr/testify/assert"
)

func TestNewUserDTO(t *testing.T) {
	user := &model.User{
		ID:    "user-123",
		Email: "alice@example.com",
		Role:  model.RoleAdmin,
	}

	dto := NewUserDTO(user)

	assert.Equal(t, "user-123", dto.ID)
	assert.Equal(t, "alice@example.com", dto.Email)
	assert.Equal(t, "admin", dto.Role)
}

func TestNewTokenPairDTO(t *testing.T) {
	pair := &model.TokenPair{
		AccessToken:  "access-token-value",
		RefreshToken: "refresh-token-value",
	}

	dto := NewTokenPairDTO(pair)

	assert.Equal(t, "access-token-value", dto.AccessToken)
	assert.Equal(t, "refresh-token-value", dto.RefreshToken)
}
