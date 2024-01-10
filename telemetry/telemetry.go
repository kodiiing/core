package telemetry

import (
	"context"
	ttrace "kodiiing/telemetry/trace"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

type ShutDownFunc func(context.Context) error

type Config struct {
	ServiceName string

	GrpcExporterEndpoint string
	HttpExporterEndpoint string
}

type telemetryProvider struct {
	serviceName          string
	grpcExporterEndpoint string
	httpExporterEndpoint string
}

func NewTelemetryProvider(cfg Config) *telemetryProvider {
	return &telemetryProvider{
		serviceName:          cfg.ServiceName,
		grpcExporterEndpoint: cfg.GrpcExporterEndpoint,
		httpExporterEndpoint: cfg.HttpExporterEndpoint,
	}
}

func (t *telemetryProvider) Run(ctx context.Context) (shutDownFuncs []ShutDownFunc, err error) {
	// resource
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(t.serviceName),
	)

	// propagator
	propagator := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)

	// create trace
	trace := ttrace.New(ttrace.Config{
		ServiceName:          t.serviceName,
		GrpcExporterEndpoint: t.grpcExporterEndpoint,
		HttpExporterEndpoint: t.httpExporterEndpoint,
	}).WithResource(res).WithGrpcExporter().WithHttpExporter()

	traceProvider, err := trace.CreateTraceProvider(ctx)
	if err != nil {
		return
	}

	otel.SetTracerProvider(traceProvider)
	otel.SetTextMapPropagator(propagator)

	shutDownFuncs = append(shutDownFuncs, traceProvider.Shutdown)

	return
}
