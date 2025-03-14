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

	exporterFuncs []func(context.Context) (trace.SpanExporter, error)
}

func New(cfg Config) *Trace {
	return &Trace{
		serviceName:          cfg.ServiceName,
		grpcExporterEndpoint: cfg.GrpcExporterEndpoint,
		httpExporterEndpoint: cfg.HttpExporterEndpoint,
	}
}

func (t *Trace) WithGrpcExporter() *Trace {
	t.exporterFuncs = append(t.exporterFuncs, func(ctx context.Context) (trace.SpanExporter, error) {
		exp, err := otlptracegrpc.New(ctx,
			otlptracegrpc.WithEndpoint(t.grpcExporterEndpoint),
			// TODO: remove this if needed
			otlptracegrpc.WithInsecure(),
		)
		if err != nil {
			return nil, err
		}

		return exp, nil
	})

	return t
}

func (t *Trace) WithHttpExporter() *Trace {
	t.exporterFuncs = append(t.exporterFuncs, func(ctx context.Context) (trace.SpanExporter, error) {
		exp, err := otlptracehttp.New(ctx,
			otlptracehttp.WithEndpoint(t.httpExporterEndpoint),
			// TODO: remove this if needed
			otlptracehttp.WithInsecure(),
		)
		if err != nil {
			return nil, err
		}

		return exp, nil
	})

	return t
}

func (t *Trace) WithResource(res *resource.Resource) *Trace {
	t.resource = res
	return t
}

func (t *Trace) CreateTraceProvider(ctx context.Context) (*sdktrace.TracerProvider, error) {
	if t.resource == nil {
		t.resource = newDefaultResource(t.serviceName)
	}

	opts := []sdktrace.TracerProviderOption{
		sdktrace.WithResource(t.resource),
	}

	for _, exporterFunc := range t.exporterFuncs {
		exp, err := exporterFunc(ctx)
		if err != nil {
			return nil, err
		}

		opts = append(opts, sdktrace.WithBatcher(exp))
	}

	return sdktrace.NewTracerProvider(opts...), nil
}

func newDefaultResource(serviceName string) *resource.Resource {
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(serviceName),
	)

	return res
}
