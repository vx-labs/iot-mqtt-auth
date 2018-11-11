package main

import (
	"log"
	"net"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/sirupsen/logrus"
	"github.com/vx-labs/iot-mqtt-auth/identity"
	"github.com/vx-labs/iot-mqtt-auth/tracing"
	"github.com/vx-labs/iot-mqtt-auth/types"
	"github.com/vx-labs/iot-mqtt-config"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type Authenticator struct {
	providers []identity.Provider
}

var ErrInvalidCredentials error = status.Error(codes.PermissionDenied, "invalid credentials")

func (a *Authenticator) Authenticate(ctx context.Context, in *types.AuthenticateRequest) (*types.AuthenticateReply, error) {
	for _, p := range a.providers {
		if p.CanHandle(in.Protocol, in.Transport) {
			identity, err := p.Authenticate(in.Protocol, in.Transport)
			if err == nil {
				logrus.Infof("identity validated by %s from %s: user is %s, scoped to tenant %s", identity.Provider, in.Transport.RemoteAddress, identity.ID, identity.Tenant)
				return &types.AuthenticateReply{
					Id:     identity.ID,
					Tenant: identity.Tenant,
					Token:  "",
				}, nil
			}
			logrus.Infof("authentication failed from %s (provider %s)", in.Transport.RemoteAddress, identity.Provider)

		}
	}
	logrus.Infof("refused authentication from %s: no provider were able to confirm remote identity", in.Transport.RemoteAddress)
	return nil, ErrInvalidCredentials
}

func main() {
	port := ":7994"
	lis, err := net.Listen("tcp", port)
	if err != nil {
		logrus.Fatalf("failed to listen: %v", err)
	}
	tracer := tracing.Instance()
	s := grpc.NewServer(
		grpc.UnaryInterceptor(
			otgrpc.OpenTracingServerInterceptor(tracer),
		),
		grpc.StreamInterceptor(
			otgrpc.OpenTracingStreamServerInterceptor(tracer),
		),
	)
	store := newAuthenticator()
	types.RegisterAuthenticationServiceServer(s, store)
	go serveHTTPHealth()
	logrus.Infof("serving authentication service on %v", port)
	if err := s.Serve(lis); err != nil {
		logrus.Fatalf("failed to serve: %v", err)
	}
}
func serveHTTPHealth() {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	log.Println(http.ListenAndServe("[::]:9000", mux))
}

func newAuthenticator() *Authenticator {
	_, vaultAPI, err := config.DefaultClients()
	if err != nil {
		panic(err)
	}
	a := &Authenticator{
		providers: []identity.Provider{
			identity.NewStaticVaultProvider(vaultAPI, "vx-psk"),
			identity.NewStaticVaultProvider(vaultAPI, "vx:psk"),
		},
	}
	return a
}
