package telemetry

import (
	"testing"
	"time"
)

func TestExporterEndpoint(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		rawEndpoint  string
		rawInsecure  string
		wantEndpoint string
		wantInsecure bool
		wantErr      bool
	}{
		{name: "http scheme is insecure", rawEndpoint: "http://collector:4317", wantEndpoint: "collector:4317", wantInsecure: true},
		{name: "https scheme uses tls", rawEndpoint: "https://collector:4317", wantEndpoint: "collector:4317", wantInsecure: false},
		{name: "endpoint without scheme", rawEndpoint: " collector:4317 ", wantEndpoint: "collector:4317", wantInsecure: true},
		{name: "env disables insecure", rawEndpoint: "collector:4317", rawInsecure: "false", wantEndpoint: "collector:4317", wantInsecure: false},
		{name: "env enables insecure", rawEndpoint: "https://collector:4317", rawInsecure: "true", wantEndpoint: "collector:4317", wantInsecure: true},
		{name: "invalid env value", rawInsecure: "yes", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			endpoint, insecure, err := exporterEndpoint(tt.rawEndpoint, tt.rawInsecure)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if endpoint != tt.wantEndpoint {
				t.Errorf("endpoint = %q, want %q", endpoint, tt.wantEndpoint)
			}
			if insecure != tt.wantInsecure {
				t.Errorf("insecure = %t, want %t", insecure, tt.wantInsecure)
			}
		})
	}
}

func TestEnvironment(t *testing.T) {
	t.Run("defaults when unset", func(t *testing.T) {
		t.Setenv("ENVIRONMENT", "")
		if got := environment(); got != defaultEnvironment {
			t.Errorf("environment() = %q, want %q", got, defaultEnvironment)
		}
	})

	t.Run("reads ENVIRONMENT", func(t *testing.T) {
		t.Setenv("ENVIRONMENT", "production")
		if got := environment(); got != "production" {
			t.Errorf("environment() = %q, want %q", got, "production")
		}
	})
}

func TestExportInterval(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want time.Duration
	}{
		{name: "default when unset", raw: "", want: defaultExportInterval},
		{name: "invalid value falls back", raw: "abc", want: defaultExportInterval},
		{name: "non-positive value falls back", raw: "0", want: defaultExportInterval},
		{name: "reads milliseconds", raw: "5000", want: 5 * time.Second},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("OTEL_METRIC_EXPORT_INTERVAL", tt.raw)
			if got := exportInterval(); got != tt.want {
				t.Errorf("exportInterval() = %v, want %v", got, tt.want)
			}
		})
	}
}
