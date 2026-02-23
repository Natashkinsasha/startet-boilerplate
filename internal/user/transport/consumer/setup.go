package consumer

import (
	"context"
	"log/slog"

	amqp091 "github.com/rabbitmq/amqp091-go"
	"golang.org/x/sync/errgroup"

	pkgamqp "starter-boilerplate/pkg/amqp"
)

type Runner func(ctx context.Context) error

func SetupConsumers(conn *amqp091.Connection) Runner {
	consumers := []*pkgamqp.Consumer{
		newExampleConsumer(),
	}

	return func(ctx context.Context) error {
		if conn == nil {
			slog.Warn("standalone mode: skipping amqp consumers")
			<-ctx.Done()
			return nil
		}

		slog.Info("amqp consumers started")

		g, ctx := errgroup.WithContext(ctx)

		for _, c := range consumers {
			g.Go(func() error {
				ch, err := conn.Channel()
				if err != nil {
					return err
				}
				defer ch.Close()

				return c.Run(ctx, ch)
			})
		}

		err := g.Wait()
		slog.Info("amqp consumers stopped")
		return err
	}
}
