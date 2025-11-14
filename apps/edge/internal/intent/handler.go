package intent

import (
	"blink/api/proto/pb"
	"net/http"

	"github.com/gin-gonic/gin"
)

type handler struct {
	recallClient pb.BlinkServiceClient
}

func (h *handler) EmitBlinkIntention(c *gin.Context) {
	rep, err := h.recallClient.BlinkEvaluation(c, &pb.BlinkEvaluationRequest{
		Blinker: c.ClientIP(),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":             rep.GetStatus().String(),
		"cooldown_remaining": rep.GetCooldownRemaining(),
	})
}
