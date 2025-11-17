package blink

import (
	"blink/lib/lock"
	"blink/lib/queue"
	"blink/lib/telemetry"
	"context"

	"github.com/jmoiron/sqlx"
)

func Register(log *telemetry.Logger, db *sqlx.DB, pubsub *queue.RabbitMQPubSub, queueName string, lock *lock.RedisLock) error {
	ctx := context.Background()

	txManager := newTracerBlinkTransactionManager(db)

	service := newService(txManager, lock)

	handler := handler{
		log:     log,
		service: service,
	}

	return pubsub.Subscribe(ctx, queueName, handler.onBlinkEvaluated)
}
