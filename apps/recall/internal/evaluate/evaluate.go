package evaluate

import (
	"blink/api/proto/pb"

	"google.golang.org/grpc"
)

func Scaffold(conn *grpc.Server) {
	h := &handler{}
	pb.RegisterBlinkServiceServer(conn, h)
}
