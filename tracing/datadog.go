package tracing

import (
	"os"

	"github.com/DataDog/dd-trace-go/tracer"
)

var tracerInstance *tracer.Tracer

func init() {
	transport := tracer.NewTransport(os.Getenv("NOMAD_IP_health"), "8126")
	tracerInstance = tracer.NewTracerTransport(transport)
}

func Instance() *tracer.Tracer {
	return tracerInstance
}
