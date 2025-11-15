package commit

import (
	"blink/lib/queue"
	"context"

	"github.com/jmoiron/sqlx"
)

func Scaffold(db *sqlx.DB, pubsub *queue.RabbitMQPubSub, queueName string) error {
	ctx := context.Background()

	handler := handler{
		db: db,
	}

	return pubsub.Subscribe(ctx, queueName, handler.OnBlinkEvaluated)
}
