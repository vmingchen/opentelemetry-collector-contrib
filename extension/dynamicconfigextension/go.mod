module github.com/vmingchen/opentelemetry-collector-contrib/extension/dynamicconfigextension

go 1.14

require (
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/dynamicconfigextension v0.0.0
	github.com/stretchr/testify v1.5.1
	github.com/vmingchen/opentelemetry-proto v0.0.0-00010101000000-000000000000
	go.opentelemetry.io/collector v0.3.1-0.20200609132241-685777fc1985
	go.uber.org/zap v1.10.0
	google.golang.org/grpc v1.29.1
)

replace github.com/vmingchen/opentelemetry-proto => ../../../opentelemetry-proto

replace github.com/open-telemetry/opentelemetry-collector-contrib/extension/dynamicconfigextension => ./
