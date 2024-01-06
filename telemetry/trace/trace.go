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

type Config struct {
	ServiceName string

	GrpcExporterEndpoint string
	HttpExporterEndpoint string
}

type Trace struct {
	serviceName string

	grpcExporterEndpoint string
	httpExporterEndpoint string

	resource *resource.Resource

	exporter []trace.SpanExporter
}

func New(cfg Config) *Trace {
	return &Trace{
		serviceName:          cfg.ServiceName,
		grpcExporterEndpoint: cfg.GrpcExporterEndpoint,
		httpExporterEndpoint: cfg.HttpExporterEndpoint,
	}
}

func (t *Trace) WithGrpcExporter(ctx context.Context) (*Trace, error) {
	exp, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(t.grpcExporterEndpoint),
	)
	if err != nil {
		return t, err
	}

	t.exporter = append(t.exporter, exp)

	return t, nil
}

func (t *Trace) WithHttpExporter(ctx context.Context) (*Trace, error) {
	exp, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(t.httpExporterEndpoint),
	)
	if err != nil {
		return t, err
	}

	t.exporter = append(t.exporter, exp)

	return t, nil
}

func (t *Trace) WithResource(res *resource.Resource) *Trace {
	t.resource = res
	return t
}

func (t *Trace) CreateTraceProvider() *sdktrace.TracerProvider {
	if t.resource == nil {
		t.resource = newDefaultResource(t.serviceName)
	}

	opts := []sdktrace.TracerProviderOption{
		sdktrace.WithResource(t.resource),
	}

	for _, exp := range t.exporter {
		opts = append(opts, sdktrace.WithBatcher(exp))
	}

	return sdktrace.NewTracerProvider(opts...)
}

func newDefaultResource(serviceName string) *resource.Resource {
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(serviceName),
	)

	return res
}
