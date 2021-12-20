module github.com/SigNoz/sample-golang-app

go 1.14

require (
	github.com/gin-gonic/gin v1.7.7
	github.com/jinzhu/gorm v1.9.12
	github.com/lib/pq v1.2.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin v0.28.0
	go.opentelemetry.io/otel v1.3.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.3.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.3.0
	go.opentelemetry.io/otel/sdk v1.3.0
	google.golang.org/grpc v1.42.0
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
)
