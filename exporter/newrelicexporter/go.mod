module github.com/open-telemetry/opentelemetry-collector-contrib/exporter/newrelicexporter

go 1.14

require (
	github.com/census-instrumentation/opencensus-proto v0.2.1
	github.com/golang/protobuf v1.3.5
	github.com/newrelic/newrelic-telemetry-sdk-go v0.2.0
	github.com/shirou/gopsutil v2.20.4+incompatible // indirect
	github.com/stretchr/testify v1.5.1
	go.opencensus.io v0.22.3
	go.opentelemetry.io/collector v0.5.0
	go.uber.org/zap v1.13.0
)
