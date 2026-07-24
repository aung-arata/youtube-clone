package telemetry

import (
	"context"
	"fmt"
	"log"
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

	// ===== BOILERPLATE =====
	// Separate context just for dialing - short timeout is fine
	// Do not reuse this ctx for exporter or resource
	dialCtx, dialCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer dialCancel()

	// ===== BOILERPLATE =====
	// gRPC connection to collector - only address changes per environment
	// ===== PROJECT SPECIFIC =====
	// "otel-collector:4317" - change this if collector has different host/port
	conn, err := grpc.DialContext(
		dialCtx,
		"otel-collector:4317",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to otel collector: %w", err)
	}

	// ===== BOILERPLATE =====
	// Exporter uses background context - not tied to dial timeout
	exporter, err := otlptracegrpc.New(context.Background(),
		otlptracegrpc.WithGRPCConn(conn),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create otlp exporter: %w", err)
	}

	// ===== BOILERPLATE =====
	// Resource uses background context too
	// ===== PROJECT SPECIFIC =====
	// serviceName is what appears in Jaeger UI dropdown
	// pass it from main() so each service has its own name
	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// ===== BOILERPLATE =====
	// TracerProvider wires everything together
	// ===== PROJECT SPECIFIC =====
	// AlwaysSample = 100% sampling, fine for dev
	// In production change to sdktrace.TraceIDRatioBased(0.1) for 10%
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	// ===== BOILERPLATE =====
	// Register globally - do not change these two blocks
	otel.SetTracerProvider(tp)

	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	// ===== BOILERPLATE =====
	// Confirm tracer initialized - helpful for debugging startup
	log.Printf("Tracer initialized successfully")

	return tp.Shutdown, nil
}
