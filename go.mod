module github.com/SigNoz/sample-golang-app

go 1.14

require (
	github.com/gin-gonic/gin v1.7.7
	github.com/mattn/go-sqlite3 v2.0.1+incompatible // indirect
	github.com/uptrace/opentelemetry-go-extra/otelgorm v0.1.12
	go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin v0.31.0
	go.opentelemetry.io/otel v1.6.3
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.6.3
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.6.3
	go.opentelemetry.io/otel/sdk v1.6.3
	go.opentelemetry.io/otel/trace v1.6.3 // indirect
	google.golang.org/grpc v1.46.0
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gorm.io/driver/sqlite v1.3.2
	gorm.io/gorm v1.23.5
)
