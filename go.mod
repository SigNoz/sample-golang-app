module github.com/rahmanfadhil/gin-bookstore

go 1.14

require (
	github.com/gin-gonic/gin v1.6.3
	github.com/jdefrank/otgorm v0.0.0-20200219012437-bfb8d99ee644
	github.com/jinzhu/gorm v1.9.12
	go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin v0.17.0
	go.opentelemetry.io/contrib/propagators v0.17.0
	go.opentelemetry.io/otel v0.17.0
	go.opentelemetry.io/otel/exporters/otlp v0.16.0
	go.opentelemetry.io/otel/sdk v0.16.0
	google.golang.org/grpc v1.35.0
)
