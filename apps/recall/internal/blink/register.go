package blink

import (
	"blink/api/proto/pb"
	"blink/lib/queue"

	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
)

func Register(conn *grpc.Server, db *sqlx.DB, pubsub *queue.RabbitMQPubSub, queueName string) {
	roTracerRepo := newTracerRepository(db)

	service := newService(roTracerRepo, pubsub, queueName)

	handler := &handler{
		service: service,
	}

	pb.RegisterEvaluationServiceServer(conn, handler)
}
