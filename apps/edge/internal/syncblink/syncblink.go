package syncblink

import (
	"github.com/gin-gonic/gin"
)

func Scaffold(r *gin.Engine) {
	r.POST("/blink", handler())
}
