package consumer

import (
	"context"
	"log/slog"

	"starter-boilerplate/internal/user/domain"
	"starter-boilerplate/internal/user/domain/model"
	"starter-boilerplate/internal/user/domain/repository"
	pkgamqp "starter-boilerplate/pkg/amqp"
)

type ProfileCreatedConsumer struct {
	repo repository.ProfileRepository
}

func NewProfileCreatedConsumer(repo repository.ProfileRepository) *ProfileCreatedConsumer {
	return &ProfileCreatedConsumer{repo: repo}
}

func (c *ProfileCreatedConsumer) Register(g *pkgamqp.ConsumerGroup) {
	pkgamqp.AddConsumer(g,
		pkgamqp.ConsumerConfig{
			Queue:                "user.profile.created",
			Exchange:             "events",
			RoutingKey:           "user.created",
			DeadLetterExchange:   "dlx",
			DeadLetterRoutingKey: "user.profile.created",
		},
		c.Consume,
	)
}

func (c *ProfileCreatedConsumer) Consume(ctx context.Context, evt *domain.UserCreatedEvent, _ pkgamqp.DeliveryMeta) error {
	slog.Info("creating profile for new user", slog.String("user_id", evt.UserID))

	profile := &model.Profile{
		UserID:  evt.UserID,
		Numbers: map[string]float64{},
		Strings: map[string]string{},
	}

	return c.repo.Upsert(ctx, profile)
}
