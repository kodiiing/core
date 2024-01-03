package telemetry

import (
	"context"
	ttrace "kodiiing/telemetry/trace"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type ShutDownFunc func(context.Context) error

type Config struct {
	ServiceName string

	GrpcEndpoint string
}

type telemetryProvider struct {
	serviceName  string
	grpcEndpoint string
}

func NewTelemetryProvider(cfg Config) *telemetryProvider {
	return &telemetryProvider{
		serviceName:  cfg.ServiceName,
		grpcEndpoint: cfg.GrpcEndpoint,
	}
}

func (t *telemetryProvider) Run(ctx context.Context) (shutDownFuncs []ShutDownFunc, err error) {
	// propagator
	propagator := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)

	// create trace
	trace := ttrace.New(ttrace.Config{
		ServiceName:  t.serviceName,
		GrpcEndpoint: t.grpcEndpoint,
	})

	trace, err = trace.WithGrpcExporter(ctx)
	if err != nil {
		return nil, err
	}

	traceProvider := trace.CreateTraceProvider()

	otel.SetTracerProvider(traceProvider)
	otel.SetTextMapPropagator(propagator)

	shutDownFuncs = append(shutDownFuncs, traceProvider.Shutdown)

	return
}
