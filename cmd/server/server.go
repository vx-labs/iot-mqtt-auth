package main

import (
	"github.com/vx-labs/iot-mqtt-auth/types"
	"google.golang.org/grpc/reflection"
	"net"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"golang.org/x/net/context"
	"os"
)

type Authenticator struct {
	logger *logrus.Entry
}

func (a *Authenticator) Authenticate(ctx context.Context, in *types.AuthenticateRequest) (*types.AuthenticateReply, error) {
	a.logger.Infof("authentication request from %s", in.Transport.RemoteAddress)
	isTransportCompliant := in.Transport.Ensure(
		types.MustBeEncrypted(),
	)
	isProtocolCompliant := in.Protocol.Ensure(
		types.MustUseStaticSharedKey(os.Getenv("PSK")),
	)
	success := isProtocolCompliant && isTransportCompliant
	return &types.AuthenticateReply{Success: success, Tenant: "_default"}, nil
}

func main() {
	port := ":7994"
	lis, err := net.Listen("tcp", port)
	if err != nil {
		logrus.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	store := &Authenticator{
		logger: logrus.New().WithField("source", "service"),
	}
	types.RegisterAuthenticationServiceServer(s, store)
	reflection.Register(s)
	logrus.Infof("serving session store on %v", port)
	if err := s.Serve(lis); err != nil {
		logrus.Fatalf("failed to serve: %v", err)
	}
}
