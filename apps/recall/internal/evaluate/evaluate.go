package evaluate

import (
	"blink/api/proto/pb"

	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
)

func Scaffold(conn *grpc.Server, db *sqlx.DB) {
	h := &handler{
		db: db,
	}

	pb.RegisterBlinkServiceServer(conn, h)
}
