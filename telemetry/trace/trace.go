package trace

import (
	"context"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

type ShutdownFunc func(context.Context) error

type Config struct {
	ServiceName string

	GrpcExporterEndpoint string
	HttpExporterEndpoint string
}

type Trace struct {
	serviceName string

	grpcExporterEndpoint string
	httpExporterEndpoint string

	exporter []trace.SpanExporter
}

func New(cfg Config) Trace {
	return Trace{
		serviceName:          cfg.ServiceName,
		grpcExporterEndpoint: cfg.GrpcExporterEndpoint,
		httpExporterEndpoint: cfg.HttpExporterEndpoint,
	}
}

func (t *Trace) WithGrpcExporter(ctx context.Context) (Trace, error) {
	exp, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(t.grpcExporterEndpoint),
	)
	if err != nil {
		return Trace{}, err
	}

	t.exporter = append(t.exporter, exp)

	return *t, nil
}

func (t *Trace) WithHttpExporter(ctx context.Context) (Trace, error) {
	exp, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(t.httpExporterEndpoint),
	)
	if err != nil {
		return Trace{}, err
	}

	t.exporter = append(t.exporter, exp)

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

	for _, exp := range t.exporter {
		opts = append(opts, sdktrace.WithBatcher(exp))
	}

	return sdktrace.NewTracerProvider(opts...)
}
