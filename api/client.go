package api

import (
	"context"
	"fmt"
	"io"

	"github.com/vx-labs/iot-mqtt-auth/types"
	"google.golang.org/grpc"
)

type Client struct {
	conn io.Closer
	api  types.AuthenticationServiceClient
}

type AuthResponse struct {
	ID     string
	Tenant string
	Token  string
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

func (c *Client) Authenticate(ctx context.Context, a ...AuthOpt) (AuthResponse, error) {
	opts := getOpts(a)
	response, err := c.api.Authenticate(ctx, &types.AuthenticateRequest{
		Transport: opts.TransportContext,
		Protocol:  opts.ProtocolContext,
	})
	if err != nil {
		return AuthResponse{}, fmt.Errorf("error ocurred when talking to authentication service: %v", err)
	}
	return AuthResponse{
		ID:     response.Id,
		Tenant: response.Tenant,
		Token:  response.Token,
	}, nil
}
