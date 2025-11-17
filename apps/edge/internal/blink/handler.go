package blink

import (
	"blink/api/proto/pb"
	"blink/lib/core"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type handler struct {
	service          *service
	evaluationClient pb.EvaluationServiceClient
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

	if cd, err := h.service.emitBlinkIntent(c.Request.Context(), req.Nickname); err != nil {
		if errors.Is(err, ErrBlinkOnCooldown) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": err.Error(),
				"data": map[string]string{
					"cooldown_remaining": fmt.Sprintf("%.2f", cd),
				},
			})
			return
		}

		if errors.Is(err, ErrProcessingBlink) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": err.Error(),
			})
			return
		}

		var notFoundErr *core.NotFoundError
		if errors.As(err, &notFoundErr) {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": notFoundErr.Error(),
			})
			return
		}

		var unknownClientErr *core.UnknowClientError
		if errors.As(err, &unknownClientErr) {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": unknownClientErr.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message": "blinking",
	})
}
