//go:build unit

package outbox

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type testEvent struct {
	Name string `json:"name"`
}

func (e testEvent) EventName() string { return "test.event" }

type taggableEvent struct {
	Name string `json:"name"`
}

func (e taggableEvent) EventName() string { return "tagged.event" }
func (e taggableEvent) Tags() []string    { return []string{"user", "profile"} }

func TestOutboxBus_Publish(t *testing.T) {
	repo := new(mockRepository)
	bus := &OutboxBus{outboxRepo: repo}

	evt := testEvent{Name: "hello"}

	repo.On("Insert", mock.Anything, mock.AnythingOfType("*outbox.Entry")).
		Run(func(args mock.Arguments) {
			entry := args.Get(1).(*Entry)

			assert.Equal(t, "test.event", entry.EventName)
			assert.NotZero(t, entry.CreatedAt)
			assert.Empty(t, entry.Headers)

			var payload map[string]any
			require.NoError(t, json.Unmarshal(entry.Payload, &payload))
			assert.Equal(t, "hello", payload["name"])
		}).
		Return(nil)

	err := bus.Publish(context.Background(), evt)

	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestOutboxBus_Publish_WithTags(t *testing.T) {
	repo := new(mockRepository)
	bus := &OutboxBus{outboxRepo: repo}

	evt := taggableEvent{Name: "tagged"}

	repo.On("Insert", mock.Anything, mock.AnythingOfType("*outbox.Entry")).
		Run(func(args mock.Arguments) {
			entry := args.Get(1).(*Entry)

			assert.Equal(t, "tagged.event", entry.EventName)
			assert.Equal(t, true, entry.Headers["tag.user"])
			assert.Equal(t, true, entry.Headers["tag.profile"])
			assert.Len(t, entry.Headers, 2)
		}).
		Return(nil)

	err := bus.Publish(context.Background(), evt)

	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestOutboxBus_Publish_RepoError(t *testing.T) {
	repo := new(mockRepository)
	bus := &OutboxBus{outboxRepo: repo}

	evt := testEvent{Name: "fail"}
	repoErr := errors.New("db connection lost")

	repo.On("Insert", mock.Anything, mock.AnythingOfType("*outbox.Entry")).
		Return(repoErr)

	err := bus.Publish(context.Background(), evt)

	assert.ErrorIs(t, err, repoErr)
	repo.AssertExpectations(t)
}
