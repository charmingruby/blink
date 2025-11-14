package grpcx

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	Conn *grpc.ClientConn
}

func NewClient(targetAddr string) (*Client, error) {
	cl, err := grpc.NewClient(targetAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &Client{
		Conn: cl,
	}, nil
}

func (c *Client) Close() error {
	return c.Conn.Close()
}
