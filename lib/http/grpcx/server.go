package grpcx

import (
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	Conn *grpc.Server
	addr string
}

func NewServer(addr string) *Server {
	srv := grpc.NewServer()

	reflection.Register(srv)

	return &Server{
		Conn: srv,
		addr: addr,
	}
}

func (s *Server) Start() error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	return s.Conn.Serve(lis)
}

func (s *Server) Stop() {
	s.Conn.GracefulStop()
}
