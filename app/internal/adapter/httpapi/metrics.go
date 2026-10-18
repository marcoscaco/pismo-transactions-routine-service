package httpapi

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

const (
	requestsInstrumentName = "http_server_requests_total"
	requestsInstrumentDesc = "Total de requisições HTTP recebidas"
	durationInstrumentName = "http_server_request_duration_seconds"
	durationInstrumentDesc = "Duração das requisições HTTP em segundos"
)

// Metrics holds the HTTP instruments shared by every request.
type Metrics struct {
	requests metric.Int64Counter
	duration metric.Float64Histogram
}

// NewMetrics builds the HTTP instruments from the given meter.
func NewMetrics(meter metric.Meter) (*Metrics, error) {
	requests, err := meter.Int64Counter(
		requestsInstrumentName,
		metric.WithDescription(requestsInstrumentDesc),
	)
	if err != nil {
		return nil, fmt.Errorf("creating %s instrument: %w", requestsInstrumentName, err)
	}

	duration, err := meter.Float64Histogram(
		durationInstrumentName,
		metric.WithDescription(durationInstrumentDesc),
	)
	if err != nil {
		return nil, fmt.Errorf("creating %s instrument: %w", durationInstrumentName, err)
	}

	return &Metrics{requests: requests, duration: duration}, nil
}

// Middleware records request count and latency using the pre-built instruments.
func (m *Metrics) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		attrs := metric.WithAttributes(
			attribute.String("http.request.method", c.Request.Method),
			attribute.String("http.route", c.FullPath()),
			attribute.Int("http.response.status_code", c.Writer.Status()),
		)

		m.requests.Add(c.Request.Context(), 1, attrs)
		m.duration.Record(c.Request.Context(), time.Since(start).Seconds(), attrs)
	}
}
