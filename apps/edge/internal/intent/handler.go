package intent

import (
	"blink/api/proto/pb"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type handler struct {
	recallClient pb.BlinkServiceClient
}

func (h *handler) EmitBlinkIntention(c *gin.Context) {
	rep, err := h.recallClient.BlinkEvaluation(c, &pb.BlinkEvaluationRequest{
		Ip: "dummy",
	})

	sts, ok := status.FromError(err)

	if !ok {
		println(sts.Code().String())

		switch sts.Code() {
		case codes.Internal:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		case codes.NotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":             rep.GetStatus().String(),
		"cooldown_remaining": rep.GetRemainingCooldown(),
	})
}
