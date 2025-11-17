package blink

import (
	"blink/api/proto/pb"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func Register(r *gin.Engine, grpcConn *grpc.ClientConn) {
	evaluationClient := pb.NewEvaluationServiceClient(grpcConn)

	service := newService(evaluationClient)

	handler := handler{
		service:          service,
		evaluationClient: evaluationClient,
	}

	r.POST("/api/blinks/intent", handler.emitBlinkIntention)
}
