package trace

import (
	"context"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

type ShutdownFunc func(context.Context) error

type Config struct {
	ServiceName string

	GrpcEndpoint string
}

type Trace struct {
	serviceName string

	grpcEndpoint string

	exporter trace.SpanExporter
}

func New(cfg Config) Trace {
	return Trace{
		serviceName:  cfg.ServiceName,
		grpcEndpoint: cfg.GrpcEndpoint,
	}
}

func (t *Trace) WithGrpcExporter(ctx context.Context) (Trace, error) {
	exp, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(t.grpcEndpoint),
	)
	if err != nil {
		return Trace{}, err
	}

	t.exporter = exp

	return *t, nil
}

func (t *Trace) CreateTraceProvider() *sdktrace.TracerProvider {
	r := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(t.serviceName),
	)

	opts := []sdktrace.TracerProviderOption{
		sdktrace.WithResource(r),
	}
	if t.exporter != nil {
		opts = append(opts, sdktrace.WithBatcher(t.exporter))
	}

	return sdktrace.NewTracerProvider(opts...)
}
