package evaluate

import (
	"blink/api/proto/pb"
	"blink/lib/queue"

	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
)

func Scaffold(conn *grpc.Server, db *sqlx.DB, pubsub *queue.RabbitMQPubSub, queueName string) {
	roTracerRepo := newTracerRepo(db)

	handler := &handler{
		repo:      roTracerRepo,
		pubsub:    pubsub,
		queueName: queueName,
	}

	pb.RegisterBlinkServiceServer(conn, handler)
}
