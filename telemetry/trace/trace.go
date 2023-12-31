package trace

import (
	"context"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

func NewTraceProvider(ctx context.Context, exp sdktrace.SpanExporter) *sdktrace.TracerProvider {
	r := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName("Kodiiing"),
	)

	return sdktrace.NewTracerProvider(
		sdktrace.WithResource(r),
		sdktrace.WithBatcher(exp),
	)
}

func CreateGrpcExporter(ctx context.Context, endpoint string) (trace.SpanExporter, error) {
	return otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(endpoint),
	)
}
