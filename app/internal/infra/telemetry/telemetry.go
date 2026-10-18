package telemetry

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"

	"github.com/marcoscaco/pismo-transactions-routine-service/internal/platform/env"
)

const (
	defaultEndpoint       = "otel-collector:4317"
	defaultEnvironment    = "development"
	defaultExportInterval = time.Second
	initTimeout           = 5 * time.Second
	shutdownTimeout       = 3 * time.Second
)

// InitTelemetry configures the global MeterProvider and TracerProvider from the
// OTEL_* environment variables and returns a function that flushes and shuts
// both providers down.
func InitTelemetry(ctx context.Context, serviceName string) (func(), error) {
	endpoint, insecure, err := exporterEndpoint(
		env.OrDefault("OTEL_EXPORTER_OTLP_ENDPOINT", defaultEndpoint),
		os.Getenv("OTEL_EXPORTER_OTLP_INSECURE"),
	)
	if err != nil {
		return func() {}, err
	}

	log.Printf("[OTel Init] Conectando no Collector gRPC em %s (insecure=%t)...", endpoint, insecure)

	ctxInit, cancelInit := context.WithTimeout(ctx, initTimeout)
	defer cancelInit()

	res, err := newResource(ctxInit, serviceName)
	if err != nil {
		log.Printf("[OTel Error] Falha ao criar resource: %v", err)
		return func() {}, err
	}

	meterProvider, err := newMeterProvider(ctxInit, res, endpoint, insecure)
	if err != nil {
		return func() {}, err
	}

	tracerProvider, err := newTracerProvider(ctxInit, res, endpoint, insecure)
	if err != nil {
		_ = meterProvider.Shutdown(ctxInit)
		return func() {}, err
	}

	otel.SetMeterProvider(meterProvider)
	otel.SetTracerProvider(tracerProvider)
	log.Println("[OTel Success] MeterProvider e TracerProvider registrados com sucesso!")

	shutdown := func() {
		ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancelShutdown()

		err := errors.Join(
			meterProvider.Shutdown(ctxShutdown),
			tracerProvider.Shutdown(ctxShutdown),
		)
		if err != nil {
			log.Printf("[OTel Shutdown Error] %v", err)
		}
	}

	return shutdown, nil
}

func newResource(ctx context.Context, serviceName string) (*resource.Resource, error) {
	// Evita conflito declarando explicitamente a SchemaURL correspondente à constante do semconv
	return resource.New(ctx,
		resource.WithSchemaURL(semconv.SchemaURL),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
			attribute.String("environment", environment()),
		),
	)
}

func newMeterProvider(
	ctx context.Context,
	res *resource.Resource,
	endpoint string,
	insecure bool,
) (*metric.MeterProvider, error) {
	opts := []otlpmetricgrpc.Option{otlpmetricgrpc.WithEndpoint(endpoint)}
	if insecure {
		opts = append(opts, otlpmetricgrpc.WithInsecure())
	}

	exporter, err := otlpmetricgrpc.New(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("criando exportador de métricas: %w", err)
	}

	reader := metric.NewPeriodicReader(exporter, metric.WithInterval(exportInterval()))

	return metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(reader),
	), nil
}

func newTracerProvider(
	ctx context.Context,
	res *resource.Resource,
	endpoint string,
	insecure bool,
) (*sdktrace.TracerProvider, error) {
	opts := []otlptracegrpc.Option{otlptracegrpc.WithEndpoint(endpoint)}
	if insecure {
		opts = append(opts, otlptracegrpc.WithInsecure())
	}

	exporter, err := otlptracegrpc.New(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("criando exportador de traces: %w", err)
	}

	return sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(exporter),
	), nil
}

// exporterEndpoint normalizes the endpoint and resolves whether the connection
// should be insecure, honouring an explicit OTEL_EXPORTER_OTLP_INSECURE value.
func exporterEndpoint(rawEndpoint, rawInsecure string) (string, bool, error) {
	endpoint := strings.TrimSpace(rawEndpoint)
	insecure := true

	switch {
	case strings.HasPrefix(endpoint, "https://"):
		endpoint = strings.TrimPrefix(endpoint, "https://")
		insecure = false
	case strings.HasPrefix(endpoint, "http://"):
		endpoint = strings.TrimPrefix(endpoint, "http://")
	}

	if rawInsecure = strings.TrimSpace(rawInsecure); rawInsecure != "" {
		parsed, err := strconv.ParseBool(rawInsecure)
		if err != nil {
			return "", false, fmt.Errorf("parsing OTEL_EXPORTER_OTLP_INSECURE: %w", err)
		}
		insecure = parsed
	}

	return endpoint, insecure, nil
}

func environment() string {
	return env.OrDefault("ENVIRONMENT", defaultEnvironment)
}

// exportInterval reads OTEL_METRIC_EXPORT_INTERVAL (milliseconds), falling back
// to the default when unset or invalid.
func exportInterval() time.Duration {
	raw := os.Getenv("OTEL_METRIC_EXPORT_INTERVAL")
	millis, err := strconv.Atoi(raw)
	if err != nil || millis <= 0 {
		return defaultExportInterval
	}
	return time.Duration(millis) * time.Millisecond
}
