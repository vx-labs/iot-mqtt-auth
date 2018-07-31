package tracing

import (
	"io"
	"log"
	"os"

	"github.com/opentracing/opentracing-go"

	datadog "github.com/DataDog/dd-trace-go/opentracing"
)

var tracerInstance opentracing.Tracer
var closer io.Closer

func init() {
	var err error
	config := datadog.NewConfiguration()
	config.AgentHostname = os.Getenv("NOMAD_IP_health")
	config.ServiceName = "mqtt-authentication"
	tracerInstance, closer, err = datadog.NewTracer(config)
	if err != nil {
		log.Fatalln(err)
	}
}

func Instance() opentracing.Tracer {
	return tracerInstance
}

func Close() error {
	return closer.Close()
}
