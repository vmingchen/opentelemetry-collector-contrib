module github.com/open-telemetry/opentelemetry-collector-contrib/extension/dynamicconfig

go 1.14

require (
	github.com/fsnotify/fsnotify v1.4.9
	github.com/open-telemetry/opentelemetry-proto v0.4.0
	github.com/spf13/viper v1.6.2
	github.com/stretchr/testify v1.6.1
	go.opentelemetry.io/collector v0.3.1-0.20200609132241-685777fc1985
	go.opentelemetry.io/contrib/exporters/metric/dynamicconfig v0.0.0-00010101000000-000000000000
	go.opentelemetry.io/otel v0.7.0
	go.opentelemetry.io/otel/exporters/otlp v0.7.0
	go.uber.org/zap v1.10.0
	google.golang.org/grpc v1.30.0
	google.golang.org/grpc/examples v0.0.0-20200630190442-3de8449f8555 // indirect
)

replace github.com/open-telemetry/opentelemetry-proto => ../../../opentelemetry-proto

replace go.opentelemetry.io/contrib => ../../../opentelemetry-go-contrib

replace go.opentelemetry.io/contrib/exporters/metric/dynamicconfig => ../../../opentelemetry-go-contrib/exporters/metric/dynamicconfig

replace github.com/open-telemetry/opentelemetry-collector-contrib/extension/dynamicconfig => ./
