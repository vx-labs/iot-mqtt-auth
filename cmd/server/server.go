package main

import (
	"log"
	"net"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/vx-labs/iot-mqtt-auth/identity"
	"github.com/vx-labs/iot-mqtt-auth/metrics"
	"github.com/vx-labs/iot-mqtt-auth/types"
	"github.com/vx-labs/iot-mqtt-config"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type Authenticator struct {
	providers []identity.Provider
}

func (a *Authenticator) Authenticate(ctx context.Context, in *types.AuthenticateRequest) (*types.AuthenticateReply, error) {
	for _, p := range a.providers {
		if p.CanHandle(in.Protocol, in.Transport) {
			identity, err := p.Authenticate(in.Protocol, in.Transport)
			if err == nil {
				logrus.Infof("identity validated by %s from %s: user is %s, scoped to tenant %s", identity.Provider, in.Transport.RemoteAddress, identity.ID, identity.Tenant)
				return &types.AuthenticateReply{
					Success: true,
					Tenant:  identity.Tenant,
				}, nil
			}
			logrus.Infof("authentication failed from %s (provider %s)", in.Transport.RemoteAddress, identity.Provider)

		}
	}
	logrus.Infof("refused authentication from %s: no provider were able to confirm remote identity", in.Transport.RemoteAddress)
	return &types.AuthenticateReply{
		Success: false,
	}, nil
}

func main() {
	port := ":7994"
	lis, err := net.Listen("tcp", port)
	if err != nil {
		logrus.Fatalf("failed to listen: %v", err)
	}
	m := metrics.NewMetricHandler()
	s := grpc.NewServer()
	store := newAuthenticator()
	types.RegisterAuthenticationServiceServer(s, store)
	go serveHTTPHealth()
	logrus.Infof("serving authentication service on %v", port)
	if err := s.Serve(lis); err != nil {
		logrus.Fatalf("failed to serve: %v", err)
	}
	m.Close()
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
