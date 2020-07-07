module github.com/vmingchen/opentelemetry-collector-contrib/extension/dynamicconfig

go 1.14

require (
	github.com/fsnotify/fsnotify v1.4.9
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/dynamicconfig v0.0.0-00010101000000-000000000000
	github.com/open-telemetry/opentelemetry-proto v0.4.0
	github.com/spf13/viper v1.6.2
	github.com/stretchr/testify v1.6.1
	go.opentelemetry.io/collector v0.5.0
	go.opentelemetry.io/contrib/exporters/metric/dynamicconfig v0.0.0-00010101000000-000000000000
	go.opentelemetry.io/otel v0.7.0
	go.opentelemetry.io/otel/exporters/otlp v0.7.0
	go.uber.org/zap v1.13.0
	google.golang.org/grpc v1.30.0
)

replace github.com/open-telemetry/opentelemetry-proto => ../../../opentelemetry-proto

replace go.opentelemetry.io/contrib => ../../../opentelemetry-go-contrib

replace go.opentelemetry.io/contrib/exporters/metric/dynamicconfig => ../../../opentelemetry-go-contrib/exporters/metric/dynamicconfig

replace github.com/open-telemetry/opentelemetry-collector-contrib/extension/dynamicconfig => ./

// force v1.29.0
replace google.golang.org/grpc => ../../../../../google.golang.org/grpc
