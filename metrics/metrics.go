package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"io"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"github.com/sirupsen/logrus"
	"net"
)

var (
	AccessDenied = prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "access_denied",
		Help:      "number of denied authentication requests",
		Namespace: "mqtt",
		Subsystem: "authentication",
	})
	AccessGranted = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name:      "access_granted",
		Help:      "number of denied authentication requests",
		Namespace: "mqtt",
		Subsystem: "authentication",
	}, []string{"method", "tenant"})
)

func init() {
	prometheus.MustRegister(AccessDenied)
	prometheus.MustRegister(AccessGranted)
}

func NewMetricHandler() io.Closer {
	http.Handle("/metrics", promhttp.Handler())
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{
		Port: 8080,
		IP:   net.IPv6zero,
	})
	if err != nil {
		panic(err)
	}
	go func() {
		logrus.Error(http.Serve(listener, nil))
	}()
	return listener
}
