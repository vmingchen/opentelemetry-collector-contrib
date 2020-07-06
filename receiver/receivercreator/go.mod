module github.com/open-telemetry/opentelemetry-collector-contrib/receiver/receivercreator

go 1.14

require (
	github.com/antonmedv/expr v1.8.4
	github.com/census-instrumentation/opencensus-proto v0.2.1
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/observer v0.0.0
	github.com/shirou/gopsutil v2.20.4+incompatible // indirect
	github.com/spf13/cast v1.3.1
	github.com/spf13/viper v1.6.2
	github.com/stretchr/testify v1.5.1
	go.opentelemetry.io/collector v0.5.0
	go.uber.org/zap v1.13.0
)

replace github.com/open-telemetry/opentelemetry-collector-contrib/extension/observer => ../../extension/observer
