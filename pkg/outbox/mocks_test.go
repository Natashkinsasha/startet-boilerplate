//go:build unit

package outbox

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type mockRepository struct {
	mock.Mock
}

func (m *mockRepository) Insert(ctx context.Context, entry *Entry) error {
	args := m.Called(ctx, entry)
	return args.Error(0)
}
