package httpapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

func TestMetricsMiddlewareRecordsRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	reader := sdkmetric.NewManualReader()
	provider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader))

	metrics, err := NewMetrics(provider.Meter("test"))
	if err != nil {
		t.Fatalf("NewMetrics() error = %v", err)
	}

	router := NewRouter(Handlers{Health: NewHealthHandler()}, metrics)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/health", nil))

	var collected metricdata.ResourceMetrics
	if err := reader.Collect(context.Background(), &collected); err != nil {
		t.Fatalf("collecting metrics: %v", err)
	}

	requests := findSum(t, collected, requestsInstrumentName)
	if len(requests.DataPoints) != 1 {
		t.Fatalf("got %d data points, want 1", len(requests.DataPoints))
	}
	if requests.DataPoints[0].Value != 1 {
		t.Errorf("requests = %v, want 1", requests.DataPoints[0].Value)
	}

	wantAttrs := map[attribute.Key]string{
		"http.request.method":       http.MethodGet,
		"http.route":                "/health",
		"http.response.status_code": "200",
	}
	for key, want := range wantAttrs {
		if got := attributeValue(requests.DataPoints[0].Attributes, key); got != want {
			t.Errorf("attribute %q = %q, want %q", key, got, want)
		}
	}

	if agg := findMetric(collected, durationInstrumentName); agg == nil {
		t.Errorf("%s not collected", durationInstrumentName)
	} else if _, ok := agg.(metricdata.Histogram[float64]); !ok {
		t.Errorf("%s has type %T, want metricdata.Histogram[float64]", durationInstrumentName, agg)
	}
}

func findMetric(rm metricdata.ResourceMetrics, name string) metricdata.Aggregation {
	for _, sm := range rm.ScopeMetrics {
		for _, m := range sm.Metrics {
			if m.Name == name {
				return m.Data
			}
		}
	}
	return nil
}

func findSum(t *testing.T, rm metricdata.ResourceMetrics, name string) metricdata.Sum[int64] {
	t.Helper()

	sum, ok := findMetric(rm, name).(metricdata.Sum[int64])
	if !ok {
		t.Fatalf("%s not found or not an int64 sum", name)
	}
	return sum
}

func attributeValue(set attribute.Set, key attribute.Key) string {
	value, ok := set.Value(key)
	if !ok {
		return ""
	}
	return value.Emit()
}
