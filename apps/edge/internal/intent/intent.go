package intent

import (
	"blink/api/proto/pb"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func Scaffold(r *gin.Engine, grpcCl *grpc.ClientConn) {
	recallCl := pb.NewBlinkServiceClient(grpcCl)

	handler := handler{
		recallClient: recallCl,
	}

	r.POST("/api/blinks/intent", handler.emitBlinkIntention)
}
