package telemetry

import (
	"context"
	tmetric "kodiiing/telemetry/metric"
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
	HttpExporterPath     string
}

type telemetryProvider struct {
	serviceName          string
	grpcExporterEndpoint string
	httpExporterEndpoint string
	httpExporterPath     string
}

func NewTelemetryProvider(cfg Config) *telemetryProvider {
	return &telemetryProvider{
		serviceName:          cfg.ServiceName,
		grpcExporterEndpoint: cfg.GrpcExporterEndpoint,
		httpExporterEndpoint: cfg.HttpExporterEndpoint,
		httpExporterPath:     cfg.HttpExporterPath,
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
	})

	trace, err = trace.WithGrpcExporter(ctx)
	if err != nil {
		return nil, err
	}
	trace, err = trace.WithHttpExporter(ctx)
	if err != nil {
		return nil, err
	}

	traceProvider := trace.CreateTraceProvider()

	otel.SetTracerProvider(traceProvider)
	otel.SetTextMapPropagator(propagator)

	// create metric
	metric := tmetric.New(tmetric.Config{
		ServiceName:          t.serviceName,
		HttpExporterEndpoint: t.httpExporterEndpoint,
		GrpcExporterEndpoint: t.grpcExporterEndpoint,
		HttpExporterPath:     t.httpExporterPath,
	}).WithResource(res).WithGrpcExporter().WithHttpExporter()

	meterProvider, err := metric.CreateMetricProvider(ctx)
	if err != nil {
		return nil, err
	}

	shutDownFuncs = append(shutDownFuncs, traceProvider.Shutdown, meterProvider.Shutdown)

	otel.SetMeterProvider(meterProvider)

	return
}
