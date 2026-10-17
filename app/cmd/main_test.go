package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// TestSetupRouter verifies that the router is properly configured with the health endpoint.
// This is a basic unit test that ensures the health check returns a 200 status with expected payload.
func TestSetupRouter(t *testing.T) {
	// Set gin to test mode to suppress logs
	gin.SetMode(gin.TestMode)

	router := SetupRouter()

	// Create a mock HTTP request for the /health endpoint
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	// Assert the response status code
	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	// Assert the response body contains expected JSON
	expectedBody := `{"status":"ok"}`
	if rec.Body.String() != expectedBody {
		t.Errorf("expected body %q, got %q", expectedBody, rec.Body.String())
	}

	// Assert the response content type is JSON
	contentType := rec.Header().Get("Content-Type")
	if contentType != "application/json; charset=utf-8" {
		t.Errorf("expected Content-Type %q, got %q", "application/json; charset=utf-8", contentType)
	}
}