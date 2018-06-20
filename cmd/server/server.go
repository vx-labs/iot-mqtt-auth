package main

import (
	"log"
	"net"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/vx-labs/iot-mqtt-auth/metrics"
	"github.com/vx-labs/iot-mqtt-auth/types"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type Authenticator struct {
	logger *logrus.Entry
}

func electTenant(left, right string) string {
	if left == "" {
		left = "_default"
	}
	if right == "" {
		right = "_default"
	}
	if left == right {
		return left
	}
	if left == "_default" {
		if right != "default" {
			return right
		}
	}
	if right == "_default" {
		if left != "default" {
			return left
		}
	}
	if left != right {
		return left
	}
	return "_default"
}

func (a *Authenticator) Authenticate(ctx context.Context, in *types.AuthenticateRequest) (*types.AuthenticateReply, error) {
	isTransportCompliant, transportTenant := in.Transport.Ensure(
		types.AlwaysAllowTransport(),
	)
	isProtocolCompliant, protocolTenant := in.Protocol.Ensure(
		types.MustUseStaticSharedKey(os.Getenv("PSK")).Or(types.MustUseStaticSharedKey(os.Getenv("PSK2"))).Or(types.MustUseDemoCredentials()),
	)
	success := isProtocolCompliant && isTransportCompliant
	if transportTenant != protocolTenant {
		a.logger.Warn("transport tenant is different from protocol tenant: %s != %s", transportTenant, protocolTenant)
		a.logger.Warn("using protocol tenant %s", protocolTenant)
	}
	tenant := electTenant(transportTenant, protocolTenant)
	if success {
		a.logger.Infof("authentication successful from %s", in.Transport.RemoteAddress)
		metrics.AccessGranted.WithLabelValues("psk", tenant).Inc()
	} else {
		a.logger.Infof("authentication denied from %s", in.Transport.RemoteAddress)
		metrics.AccessDenied.Inc()
	}
	return &types.AuthenticateReply{Success: success, Tenant: tenant}, nil
}

func main() {
	port := ":7994"
	lis, err := net.Listen("tcp", port)
	if err != nil {
		logrus.Fatalf("failed to listen: %v", err)
	}
	m := metrics.NewMetricHandler()
	s := grpc.NewServer()
	store := &Authenticator{
		logger: logrus.New().WithField("source", "service"),
	}
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
