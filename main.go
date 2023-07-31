package main

import (
	"context"
	"log"
	"os"

	"github.com/SigNoz/sample-golang-app/controllers"
	"github.com/SigNoz/sample-golang-app/models"
	"google.golang.org/grpc/credentials"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"

	metricsdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"math/rand"
	"time"

	api "go.opentelemetry.io/otel/metric"
)

var (
	serviceName  = os.Getenv("SERVICE_NAME")
	collectorURL = os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	insecure     = os.Getenv("INSECURE_MODE")
)

// Example use cases for sync counter
// - count the number of bytes received
// - count the number of requests completed
// - count the number of accounts created
// - count the number of checkpoints run
// - count the number of HTTP 5xx errors
//
// The increments should be non-negative.

func exceptionsCounter(meter api.Meter) {
	counter, err := meter.Int64Counter("exceptions", api.WithUnit("1"),
		api.WithDescription("Counts exceptions since start"),
	)
	if err != nil {
		log.Fatal(err)
	}

	for {
		// Increment the counter by 1.
		// The attributes describe the exception.
		counter.Add(context.Background(), 1, api.WithAttributes(attribute.KeyValue{
			Key: attribute.Key("exception_type"), Value: attribute.StringValue("NullPointerException"),
		}))
		time.Sleep(time.Duration(rand.Int63n(5)) * time.Millisecond)
	}
}

// Example use cases for async counter
//   - count the number of page faults
//   - CPU time, which could be reported for each thread, each process or the
//     entire system. For example "the CPU time for process
//     A running in user mode, measured in seconds".
//
// Basically, any value that is monotonically increasing and happens in the background.
// The increments should be non-negative.
func pageFaultsCounter(meter api.Meter) {
	counter, err := meter.Int64ObservableCounter(
		"page_faults",
		api.WithUnit("1"),
		api.WithDescription("Counts page faults since start"),
	)
	if err != nil {
		log.Fatal(err)
	}

	_, err = meter.RegisterCallback(
		func(_ context.Context, o api.Observer) error {
			attrSet := attribute.NewSet(attribute.String("process", "boo"))
			withAttrSet := api.WithAttributeSet(attrSet)
			o.ObserveInt64(counter, rand.Int63n(100), withAttrSet)
			return nil
		},
		counter,
	)
}

// Example use cases for Histogram
// - the request duration
// - the size of the response payload
func requestDurationHistogram(meter api.Meter) {
	histogram, err := meter.Int64Histogram(
		"http_request_duration",
		api.WithUnit("ms"),
		api.WithDescription("The HTTP request duration in milliseconds"),
	)
	if err != nil {
		log.Fatal(err)
	}

	for {
		histogram.Record(context.Background(), rand.Int63n(1000), api.WithAttributes(attribute.String("path", "/api/boo")))
		time.Sleep(time.Duration(rand.Int63n(5)) * time.Millisecond)
	}
}

// Asynchronous Gauge is an Instrument
// which reports non-additive value(s)
// (e.g. the room temperature - it makes no sense to report the
// temperature value from multiple rooms and sum them up) when the
// instrument is being observed.

// Example use cases for Async Gauge
// - the current room temperature
// - the CPU fan speed

// Note: if the values are additive (e.g. the process heap size -
// it makes sense to report the heap size from multiple processes and sum them up,
// so we get the total heap usage),
// use Asynchronous Counter or Asynchronous UpDownCounter.

func roomTemperatureGauge(meter api.Meter) {
	gauge, err := meter.Float64ObservableGauge(
		"room_temperature",
		api.WithUnit("1"),
		api.WithDescription("The room temperature in celsius"),
	)
	if err != nil {
		log.Fatal(err)
	}

	_, err = meter.RegisterCallback(
		func(_ context.Context, o api.Observer) error {
			attrSet := attribute.NewSet(attribute.String("process", "boo"))
			withAttrSet := api.WithAttributeSet(attrSet)
			o.ObserveFloat64(gauge, rand.Float64()*100, withAttrSet)
			return nil
		},
		gauge,
	)
}

// UpDownCounter is an Instrument which supports increments and decrements.
// if the value is monotonically increasing, use Counter instead.
// Example use cases for UpDownCounter
// - the number of active requests
// - the number of items in a queue

func itemsInQueueUpDownCounter(meter api.Meter) {
	counter, err := meter.Int64UpDownCounter(
		"items_in_queue",
		api.WithUnit("1"),
		api.WithDescription("The number of items in the queue"),
	)
	if err != nil {
		log.Fatal(err)
	}

	for {
		counter.Add(context.Background(), rand.Int63n(100), api.WithAttributes(attribute.String("queue", "A")))
		time.Sleep(time.Duration(rand.Int63n(5)) * time.Millisecond)
	}
}

// Asynchronous UpDownCounter is an asynchronous Instrument
// which reports additive value(s)
// (e.g. the process heap size - it makes sense to report the heap size
// from multiple processes and sum them up, so we get the total heap usage)
// when the instrument is being observed.
//
// Example use cases for Asynchronous UpDownCounter
// - the process heap size
// - the approximate number of items in a lock-free circular buffer

func processHeapSizeUpDownCounter(meter api.Meter) {
	counter, err := meter.Float64ObservableUpDownCounter(
		"process_heap_size",
		api.WithUnit("1"),
		api.WithDescription("The process heap size"),
	)
	if err != nil {
		log.Fatal(err)
	}

	_, err = meter.RegisterCallback(
		func(_ context.Context, o api.Observer) error {
			attrSet := attribute.NewSet(attribute.String("process", "boo"))
			withAttrSet := api.WithAttributeSet(attrSet)
			o.ObserveFloat64(counter, rand.Float64()*100, withAttrSet)
			return nil
		},
		counter,
	)
}

func initTracer() func(context.Context) error {

	secureOption := otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
	if len(insecure) > 0 {
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
		log.Fatal(err)
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

func initMeter() *metricsdk.MeterProvider {

	secureOption := otlpmetricgrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
	if len(insecure) > 0 {
		secureOption = otlpmetricgrpc.WithInsecure()
	}

	exporter, err := otlpmetricgrpc.New(
		context.Background(),
		secureOption,
		otlpmetricgrpc.WithEndpoint(collectorURL),
	)

	if err != nil {
		log.Fatal(err)
	}

	res, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", serviceName),
			attribute.String("library.language", "go"),
		),
	)
	if err != nil {
		log.Fatalf("Could not set resources: %v", err)
	}

	// Register the exporter with an SDK via a periodic reader.
	provider := metricsdk.NewMeterProvider(
		metricsdk.WithResource(res),
		metricsdk.WithReader(metricsdk.NewPeriodicReader(exporter)),
	)
	return provider
}

func main() {

	cleanup := initTracer()
	defer cleanup(context.Background())

	provider := initMeter()
	defer provider.Shutdown(context.Background())

	meter := provider.Meter("sample-golang-app")

	go exceptionsCounter(meter)
	go pageFaultsCounter(meter)
	go requestDurationHistogram(meter)
	go roomTemperatureGauge(meter)
	go itemsInQueueUpDownCounter(meter)
	go processHeapSizeUpDownCounter(meter)

	r := gin.Default()
	r.Use(otelgin.Middleware(serviceName))
	// Connect to database
	models.ConnectDatabase()

	// Routes
	r.GET("/books", controllers.FindBooks)
	r.GET("/books/:id", controllers.FindBook)
	r.POST("/books", controllers.CreateBook)
	r.PATCH("/books/:id", controllers.UpdateBook)
	r.DELETE("/books/:id", controllers.DeleteBook)

	// Run the server
	r.Run(":8090")
}
