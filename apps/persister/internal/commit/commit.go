package commit

import (
	"blink/lib/queue"
	"blink/lib/telemetry"
	"context"

	"github.com/jmoiron/sqlx"
)

func Scaffold(log *telemetry.Logger, db *sqlx.DB, pubsub *queue.RabbitMQPubSub, queueName string) error {
	ctx := context.Background()

	txManager := newTracerBlinkTransactionManager(db)

	handler := handler{
		log:       log,
		txManager: txManager,
	}

	return pubsub.Subscribe(ctx, queueName, handler.onBlinkEvaluated)
}
