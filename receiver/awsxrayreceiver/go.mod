module github.com/open-telemetry/opentelemetry-collector-contrib/receiver/awsxrayreceiver

go 1.14

replace github.com/open-telemetry/opentelemetry-collector-contrib/internal/common => ../../internal/common

require (
	github.com/aws/aws-sdk-go v1.37.15
	github.com/google/uuid v1.1.1
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/common v0.0.0-00010101000000-000000000000 // indirect
	github.com/stretchr/testify v1.6.1
	go.opentelemetry.io/collector v0.8.1-0.20200818152037-30c3c343c558
	go.uber.org/zap v1.15.0
)
