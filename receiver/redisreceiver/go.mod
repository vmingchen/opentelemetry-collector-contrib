module github.com/open-telemetry/opentelemetry-collector-contrib/receiver/redisreceiver

go 1.14

require (
	github.com/census-instrumentation/opencensus-proto v0.3.0
	github.com/go-redis/redis/v7 v7.4.0
	github.com/golang/protobuf v1.4.2
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/common v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.7.0
	go.opentelemetry.io/collector v0.8.1-0.20200818152037-30c3c343c558
	go.uber.org/zap v1.17.0
)

replace github.com/open-telemetry/opentelemetry-collector-contrib/internal/common => ../../internal/common
