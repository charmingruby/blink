package intent

import (
	"blink/api/proto/pb"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type handler struct {
	recallClient pb.BlinkServiceClient
}

type emitBlinkIntentionRequest struct {
	Nickname string `json:"nickname" binding:"required"`
}

func (h *handler) emitBlinkIntention(c *gin.Context) {
	var req emitBlinkIntentionRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	rep, err := h.recallClient.BlinkEvaluation(c, &pb.BlinkEvaluationRequest{
		Nickname: req.Nickname,
	})

	sts, ok := status.FromError(err)

	if !ok {
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

	res := gin.H{
		"status":             rep.GetStatus().String(),
		"cooldown_remaining": fmt.Sprintf("%.2f", rep.GetRemainingCooldown()),
	}

	if rep.RemainingCooldown != 0 {
		c.JSON(http.StatusTooManyRequests, res)
		return
	}

	c.JSON(http.StatusAccepted, res)
}
