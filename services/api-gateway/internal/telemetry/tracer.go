package telemetry

import (
	"context"
	"fmt"
	"time"

	// otel core - the API layer, defines interfaces

	// the exporter that speaks OTLP gRPC to the collector
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"

	// W3C traceparent header propagator

	// what identifies your service in Jaeger
	"go.opentelemetry.io/otel/sdk/resource"
	// the SDK layer - actual implementation of those interfaces
	// aliased to avoid conflict with runtime/trace or otel/trace
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	// semantic conventions - standard attribute names
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitTracer(serviceName string) (func(context.Context) error, error) {

	// Create a gRPC connection to the OTel Collector
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		"otel-collector:4317",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to otel collector: %w", err)
	}

	// Create the OTLP exporter using that connection
	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithGRPCConn(conn),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create otlp exporter: %w", err)
	}

	// Create a Resource
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
		),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Create the TraceProvider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	// Register TraceProvider globally
	otel.SetTracerProvider(tp)

	// Register W3C propagator globally
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	return tp.Shutdown, nil
}
