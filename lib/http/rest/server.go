package rest

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

type Server struct {
	Mux  *gin.Engine
	conn http.Server
}

func NewServer(serviceName, port string) *Server {
	addr := ":" + port

	r := gin.Default()

	r.Use(otelgin.Middleware(serviceName))

	return &Server{
		conn: http.Server{
			Handler:      r,
			Addr:         addr,
			ReadTimeout:  10 * time.Second,
			IdleTimeout:  15 * time.Second,
			WriteTimeout: 30 * time.Second,
		},
		Mux: r,
	}
}

func (s *Server) Start() error {
	if err := s.conn.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.conn.Shutdown(ctx)
}
