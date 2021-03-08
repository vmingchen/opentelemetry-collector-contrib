module github.com/open-telemetry/opentelemetry-collector-contrib/receiver/sapmreceiver

go 1.14

require (
	github.com/Azure/go-autorest/autorest/adal v0.9.0 // indirect
	github.com/gorilla/mux v1.7.4
	github.com/jaegertracing/jaeger v1.22.0
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/common v0.0.0-00010101000000-000000000000
	github.com/open-telemetry/opentelemetry-proto v0.4.0
	github.com/signalfx/sapm-proto v0.7.0
	github.com/stretchr/testify v1.7.0
	go.opencensus.io v0.23.0
	go.opentelemetry.io/collector v0.8.1-0.20200818152037-30c3c343c558
	go.uber.org/zap v1.16.0
)

replace github.com/open-telemetry/opentelemetry-collector-contrib/internal/common => ../../internal/common
