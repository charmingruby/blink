package telemetry

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.34.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type TracerConfig struct {
	ServiceName string
	Endpoint    string
}

func NewTracer(cfg TracerConfig) error {
	ctx := context.Background()

	exporter, err := newTracerExporter(ctx, cfg.Endpoint)
	if err != nil {
		return err
	}

	provider, err := newTracerProvider(ctx, exporter, cfg.ServiceName)
	if err != nil {
		return err
	}

	otel.SetTracerProvider(provider)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return nil
}

func newTracerExporter(ctx context.Context, endpoint string) (sdktrace.SpanExporter, error) {
	grpcConn, err := grpc.NewClient(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(grpcConn))
	if err != nil {
		return nil, err
	}

	return exporter, nil
}

func newTracerProvider(
	ctx context.Context,
	exporter sdktrace.SpanExporter,
	serviceName string,
) (*sdktrace.TracerProvider, error) {
	resource, err := resource.New(ctx, resource.WithAttributes(
		semconv.ServiceName(serviceName),
	))
	if err != nil {
		return nil, err
	}

	provider := sdktrace.NewTracerProvider(sdktrace.WithBatcher(exporter), sdktrace.WithResource(resource))

	return provider, nil
}
