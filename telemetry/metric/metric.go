package metric

import (
	"context"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

type Config struct {
	ServiceName string

	GrpcExporterEndpoint string
	HttpExporterEndpoint string
	HttpExporterPath     string
}

type Metric struct {
	serviceName string

	grpcExporterEndpoint string
	httpExporterEndpoint string
	httpExporterPath     string

	resource *resource.Resource

	exporterFuncs []func(context.Context) (sdkmetric.Exporter, error)
}

func New(cfg Config) *Metric {
	return &Metric{
		serviceName: cfg.ServiceName,

		grpcExporterEndpoint: cfg.GrpcExporterEndpoint,
		httpExporterEndpoint: cfg.HttpExporterEndpoint,
		httpExporterPath:     cfg.HttpExporterPath,
	}
}

func (m *Metric) WithGrpcExporter() *Metric {
	m.exporterFuncs = append(m.exporterFuncs, func(ctx context.Context) (sdkmetric.Exporter, error) {
		exp, err := otlpmetricgrpc.New(ctx,
			otlpmetricgrpc.WithEndpoint(m.grpcExporterEndpoint),
			otlpmetricgrpc.WithCompressor("gzip"),
		)
		if err != nil {
			return nil, err
		}

		return exp, nil
	})

	return m
}

func (m *Metric) WithResource(res *resource.Resource) *Metric {
	m.resource = res
	return m
}

func (m *Metric) WithHttpExporter() *Metric {
	m.exporterFuncs = append(m.exporterFuncs, func(ctx context.Context) (sdkmetric.Exporter, error) {
		exp, err := otlpmetrichttp.New(ctx,
			otlpmetrichttp.WithEndpoint(m.httpExporterEndpoint),
			otlpmetrichttp.WithURLPath(m.httpExporterPath),
			otlpmetrichttp.WithCompression(otlpmetrichttp.GzipCompression),
		)
		if err != nil {
			return nil, err
		}

		return exp, nil
	})

	return m
}

func (m *Metric) CreateMetricProvider(ctx context.Context) (*sdkmetric.MeterProvider, error) {
	if m.resource == nil {
		m.resource = newDefaultResource(m.serviceName)
	}

	opts := []sdkmetric.Option{
		sdkmetric.WithResource(m.resource),
	}

	for _, exporterFunc := range m.exporterFuncs {
		exp, err := exporterFunc(ctx)
		if err != nil {
			return nil, err
		}

		opts = append(opts, sdkmetric.WithReader(
			sdkmetric.NewPeriodicReader(exp),
		))
	}

	return sdkmetric.NewMeterProvider(opts...), nil
}

func newDefaultResource(serviceName string) *resource.Resource {
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(serviceName),
	)

	return res
}
