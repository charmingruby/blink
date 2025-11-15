package evaluate

import (
	"blink/api/proto/pb"
	"blink/lib/queue"

	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
)

func Scaffold(conn *grpc.Server, db *sqlx.DB, pubsub *queue.RabbitMQPubSub) {
	roTracerRepo := newTracerRepo(db)

	handler := &handler{
		repo:   roTracerRepo,
		pubsub: pubsub,
	}

	pb.RegisterBlinkServiceServer(conn, handler)
}
