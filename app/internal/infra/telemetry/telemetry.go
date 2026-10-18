package telemetry

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
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

// InitTelemetry configures the global MeterProvider from the OTEL_* environment
// variables and returns a function that flushes and shuts the provider down.
func InitTelemetry(ctx context.Context, serviceName string) (func(), error) {
	rawEndpoint := env.OrDefault("OTEL_EXPORTER_OTLP_ENDPOINT", defaultEndpoint)

	endpoint, insecure, err := exporterEndpoint(
		rawEndpoint,
		os.Getenv("OTEL_EXPORTER_OTLP_INSECURE"),
	)
	if err != nil {
		return func() {}, err
	}

	log.Printf("[OTel Init] Conectando no Collector gRPC em %s (insecure=%t)...", endpoint, insecure)

	exporterOpts := []otlpmetricgrpc.Option{otlpmetricgrpc.WithEndpoint(endpoint)}
	if insecure {
		exporterOpts = append(exporterOpts, otlpmetricgrpc.WithInsecure())
	}

	ctxTimeout, cancel := context.WithTimeout(ctx, initTimeout)
	defer cancel()

	metricExporter, err := otlpmetricgrpc.New(ctxTimeout, exporterOpts...)
	if err != nil {
		return func() {}, err
	}

	// Evita conflito declarando explicitamente a SchemaURL correspondente à constante do semconv
	res, err := resource.New(ctx,
		resource.WithSchemaURL(semconv.SchemaURL),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
			attribute.String("environment", environment()),
		),
	)
	if err != nil {
		log.Printf("[OTel Error] Falha ao criar resource: %v", err)
		return func() {}, err
	}

	reader := metric.NewPeriodicReader(
		metricExporter,
		metric.WithInterval(exportInterval()),
	)

	meterProvider := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(reader),
	)

	otel.SetMeterProvider(meterProvider)
	log.Println("[OTel Success] MeterProvider registrado com sucesso!")

	shutdown := func() {
		ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancelShutdown()
		if err := meterProvider.Shutdown(ctxShutdown); err != nil {
			log.Printf("[OTel Shutdown Error] %v", err)
		}
	}

	return shutdown, nil
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
