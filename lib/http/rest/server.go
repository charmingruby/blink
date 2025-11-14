package rest

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
	Mux *gin.Engine
	http.Server
}

func NewServer(port string) *Server {
	addr := ":" + port

	r := gin.Default()

	return &Server{
		Server: http.Server{
			Handler:      r,
			Addr:         addr,
			ReadTimeout:  10 * time.Second,
			IdleTimeout:  15 * time.Second,
			WriteTimeout: 30 * time.Second,
		},
		Mux: r,
	}
}
