package intent

import (
	"blink/api/proto/pb"

	"github.com/gin-gonic/gin"
)

type handler struct {
	recallClient pb.BlinkServiceClient
}

func (h *handler) EmitBlinkIntention(ctx *gin.Context) {
	h.recallClient.BlinkEvaluation(ctx, &pb.BlinkEvaluationRequest{
		Blinker: ctx.ClientIP(),
	})
}
