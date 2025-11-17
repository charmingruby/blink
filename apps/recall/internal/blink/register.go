package blink

import (
	"blink/api/proto/pb"
	"blink/lib/lock"
	"blink/lib/queue"

	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
)

func Register(conn *grpc.Server, db *sqlx.DB, pubsub *queue.RabbitMQPubSub, queueName string, lock *lock.RedisLock) {
	roTracerRepo := newTracerRepository(db)

	service := newService(roTracerRepo, pubsub, queueName, lock)

	handler := &handler{
		service: service,
	}

	pb.RegisterEvaluationServiceServer(conn, handler)
}
