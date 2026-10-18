//go:build e2e

package integration

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/marcoscaco/pismo-transactions-routine-service/internal/platform/env"
)

// healthResponse represents the expected payload returned by the /health endpoint.
type healthResponse struct {
	Status string `json:"status"`
}

// TestHealthCheckE2E performs a blackbox E2E test of the healthcheck endpoint.
// It makes real HTTP requests against an externally running service whose
// address is configured via the BASE_URL environment variable.
// The test does NOT start its own server—it is completely independent of the
// application under test, making it reusable against any environment
// (local, staging, production, etc.).
//
// Environment variable:
//   - BASE_URL: the base URL of the service (default: "http://localhost:8080")
//     Example: "https://api.example.com", "http://localhost:8080", "http://my-service:80"
func TestHealthCheckE2E(t *testing.T) {
	baseURL := env.OrDefault("BASE_URL", "http://localhost:8080")

	t.Logf("Testing healthcheck against: %s/health", baseURL)

	req, err := http.NewRequest(http.MethodGet, baseURL+"/health", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("failed to reach %s/health: %v", baseURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType != "application/json; charset=utf-8" {
		t.Errorf("expected Content-Type %q, got %q", "application/json; charset=utf-8", contentType)
	}

	var body healthResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	if body.Status != "ok" {
		t.Errorf("expected status %q, got %q", "ok", body.Status)
	}
}
