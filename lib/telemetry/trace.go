package telemetry

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.34.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var tracer trace.Tracer

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
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	tracer = otel.Tracer(cfg.ServiceName)

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

func StartSpan(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return tracer.Start(ctx, spanName, opts...)
}

func StartSpanWithTracer(ctx context.Context, tracerName, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	t := otel.Tracer(tracerName)
	return t.Start(ctx, spanName, opts...)
}

func GetSpan(ctx context.Context) trace.Span {
	return trace.SpanFromContext(ctx)
}

func GetTracer(name string) trace.Tracer {
	return otel.Tracer(name)
}

func RecordError(ctx context.Context, err error) {
	span := trace.SpanFromContext(ctx)
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
}

func SetAttributes(ctx context.Context, attrs ...attribute.KeyValue) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attrs...)
}

func SetAttribute(ctx context.Context, key, value string) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String(key, value))
}

func AddEvent(ctx context.Context, name string, attrs ...attribute.KeyValue) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent(name, trace.WithAttributes(attrs...))
}
