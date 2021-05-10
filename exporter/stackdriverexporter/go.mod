module github.com/open-telemetry/opentelemetry-collector-contrib/exporter/stackdriverexporter

go 1.14

require (
	contrib.go.opencensus.io/exporter/stackdriver v0.13.3
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace v0.2.2-0.20200728233621-2752da7eaab7
	github.com/golang/protobuf v1.5.2
	github.com/stretchr/testify v1.6.1
	go.opencensus.io v0.23.0
	go.opentelemetry.io/collector v0.8.1-0.20200818152037-30c3c343c558
	go.opentelemetry.io/otel v0.9.0
	go.uber.org/zap v1.15.0
	google.golang.org/api v0.46.0
	google.golang.org/genproto v0.0.0-20210429181445-86c259c2b4ab
	google.golang.org/grpc v1.37.0
	google.golang.org/grpc/examples v0.0.0-20200728194956-1c32b02682df // indirect
)
