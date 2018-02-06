package api

import (
	"github.com/vx-labs/iot-mqtt-auth/types"
	"google.golang.org/grpc"
	"io"
)

type Client struct {
	conn io.Closer
	api  types.AuthenticationServiceClient
}

func New(addr string) (*Client, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	c := &Client{
		conn: conn,
		api:  types.NewAuthenticationServiceClient(conn),
	}
	return c, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
