package main

import (
	"context"
	"embed"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/SigNoz/sample-golang-app/controllers"
	"github.com/SigNoz/sample-golang-app/metrics"
	"github.com/SigNoz/sample-golang-app/models"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/log/global"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc/credentials"

	"github.com/gin-gonic/gin"
)

//go:embed static
var staticFS embed.FS

// indexContent is loaded at startup so we serve from memory (avoids path/FS lookup issues).
var indexContent []byte

var (
	serviceName  = os.Getenv("SERVICE_NAME")
	collectorURL = os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	insecure     = os.Getenv("INSECURE_MODE")
)

func initTracer() func(context.Context) error {
	var secureOption otlptracegrpc.Option

	if strings.ToLower(insecure) == "false" || insecure == "0" || strings.ToLower(insecure) == "f" {
		secureOption = otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
	} else {
		secureOption = otlptracegrpc.WithInsecure()
	}

	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			secureOption,
			otlptracegrpc.WithEndpoint(collectorURL),
		),
	)
	if err != nil {
		log.Fatalf("Failed to create exporter: %v", err)
	}

	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", serviceName),
			attribute.String("library.language", "go"),
		),
	)
	if err != nil {
		log.Fatalf("Could not set resources: %v", err)
	}

	otel.SetTracerProvider(
		sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
			sdktrace.WithBatcher(exporter),
			sdktrace.WithResource(resources),
		),
	)
	return exporter.Shutdown
}

func initLogger() func(context.Context) error {
	var opts []otlploggrpc.Option
	if strings.ToLower(insecure) == "false" || insecure == "0" || strings.ToLower(insecure) == "f" {
		opts = append(opts, otlploggrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, "")))
	} else {
		opts = append(opts, otlploggrpc.WithInsecure())
	}
	opts = append(opts, otlploggrpc.WithEndpoint(collectorURL))

	exporter, err := otlploggrpc.New(context.Background(), opts...)
	if err != nil {
		log.Fatalf("Failed to create log exporter: %v", err)
	}

	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", serviceName),
			attribute.String("library.language", "go"),
		),
	)
	if err != nil {
		log.Fatalf("Could not set log resources: %v", err)
	}

	processor := sdklog.NewBatchProcessor(exporter)
	provider := sdklog.NewLoggerProvider(
		sdklog.WithProcessor(processor),
		sdklog.WithResource(resources),
	)
	global.SetLoggerProvider(provider)
	// 使用 otelslog bridge 上报到 OTLP，同时用 multiHandler 写入 stdout，这样 slog 既打控制台也上报
	otelLogger := otelslog.NewLogger(serviceName, otelslog.WithLoggerProvider(provider))
	stdoutHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{})
	slog.SetDefault(slog.New(&multiHandler{handlers: []slog.Handler{stdoutHandler, otelLogger.Handler()}}))
	return provider.Shutdown
}

// multiHandler 将每条日志转发给多个 slog.Handler（同时打 stdout 和 OTLP）
type multiHandler struct {
	handlers []slog.Handler
}

func (m *multiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, h := range m.handlers {
		if h.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (m *multiHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, h := range m.handlers {
		if err := h.Handle(ctx, r); err != nil {
			return err
		}
	}
	return nil
}

func (m *multiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	next := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		next[i] = h.WithAttrs(attrs)
	}
	return &multiHandler{handlers: next}
}

func (m *multiHandler) WithGroup(name string) slog.Handler {
	next := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		next[i] = h.WithGroup(name)
	}
	return &multiHandler{handlers: next}
}

func main() {
	// Initialize the tracer
	cleanup := initTracer()
	defer cleanup(context.Background())

	// Initialize the logger (export logs via OpenTelemetry OTLP)
	logCleanup := initLogger()
	defer logCleanup(context.Background())

	// Initialize the metrics
	provider := metrics.InitMeter()
	defer provider.Shutdown(context.Background())

	meter := provider.Meter("sample-golang-app")
	metrics.GenerateMetrics(meter)

	// Connect to database
	models.ConnectDatabase()

	// Initialize the Gin server (disable redirects to avoid 301 on API calls)
	r := gin.Default()
	// r.RedirectTrailingSlash = false
	// r.RedirectFixedPath = false
	// Only normalize empty path to "/" so root URL works without any redirect
	r.Use(func(c *gin.Context) {
		if c.Request.URL.Path == "" {
			c.Request.URL.Path = "/"
		}
		c.Next()
	})
	r.Use(otelgin.Middleware(serviceName))

	// Serve embedded frontend from memory
	r.GET("/", func(c *gin.Context) {
		// Here the index page was renamed from index.html to index.htm,
		// check this issue: https://github.com/gin-gonic/gin/issues/2654.
		c.FileFromFS("static/index.htm", http.FS(staticFS))
	})
	r.NoRoute(func(c *gin.Context) {
		c.AbortWithStatus(http.StatusNotFound)
	})

	// Initialize the routes
	r.GET("/books", controllers.FindBooks)
	r.GET("/books/:id", controllers.FindBook)
	r.POST("/books", controllers.CreateBook)
	r.PATCH("/books/:id", controllers.UpdateBook)
	r.DELETE("/books/:id", controllers.DeleteBook)

	// Run the server
	slog.Info("Starting server on port 8090")
	slog.Info("You can access the WebUI at http://localhost:8090, then you test the API endpoints")

	r.Run(":8090")
}
