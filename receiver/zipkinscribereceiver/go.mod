module github.com/open-telemetry/opentelemetry-collector-contrib/receiver/zipkinscribereceiver

go 1.14

require (
	github.com/apache/thrift v0.0.0-20161221203622-b2a4d4ae21c7
	github.com/census-instrumentation/opencensus-proto v0.2.1
	github.com/golang/protobuf v1.3.5
	github.com/jaegertracing/jaeger v1.17.0
	github.com/omnition/scribe-go v1.0.0
	github.com/stretchr/testify v1.5.1
	go.opentelemetry.io/collector v0.4.0
	go.uber.org/zap v1.10.0
)
