package blink

import (
	"blink/lib/queue"
	"blink/lib/telemetry"
	"context"

	"github.com/jmoiron/sqlx"
)

func Register(log *telemetry.Logger, db *sqlx.DB, pubsub *queue.RabbitMQPubSub, queueName string) error {
	ctx := context.Background()

	txManager := newTracerBlinkTransactionManager(db)

	service := newService(txManager)

	handler := handler{
		log:     log,
		service: service,
	}

	return pubsub.Subscribe(ctx, queueName, handler.onBlinkEvaluated)
}
