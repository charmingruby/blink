package blink

import (
	"github.com/gin-gonic/gin"
)

func handler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// sync call, in case positive, dispatch to persister
		// return to user respective response
	}
}
