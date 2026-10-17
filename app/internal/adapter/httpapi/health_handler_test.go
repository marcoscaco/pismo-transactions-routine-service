package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// TestRouterHealthCheck verifies that the router exposes the health endpoint
// with the expected status, payload and content type.
func TestRouterHealthCheck(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := NewRouter(Handlers{Health: NewHealthHandler()})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	expectedBody := `{"status":"ok"}`
	if rec.Body.String() != expectedBody {
		t.Errorf("expected body %q, got %q", expectedBody, rec.Body.String())
	}

	contentType := rec.Header().Get("Content-Type")
	if contentType != "application/json; charset=utf-8" {
		t.Errorf("expected Content-Type %q, got %q", "application/json; charset=utf-8", contentType)
	}
}
