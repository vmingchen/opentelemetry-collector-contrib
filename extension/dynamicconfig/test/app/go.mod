module github.com/vmingchen/opentelemetry-collector-contrib/extension/dynamicconfig/test/app

go 1.14

require (
	go.opentelemetry.io/contrib/exporters/metric/dynamicconfig v0.0.0-00010101000000-000000000000 // indirect
	go.opentelemetry.io/otel/exporters/otlp v0.7.0 // indirect
)

replace github.com/open-telemetry/opentelemetry-proto => ../../../../../opentelemetry-proto

replace go.opentelemetry.io/contrib => ../../../../../opentelemetry-go-contrib

replace go.opentelemetry.io/contrib/exporters/metric/dynamicconfig => ../../../../../opentelemetry-go-contrib/exporters/metric/dynamicconfig
