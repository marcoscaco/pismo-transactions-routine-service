package env

import "testing"

func TestOrDefault(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		fallback string
		want     string
	}{
		{name: "uses the value when set", value: "http://api:8080", fallback: "http://localhost:8080", want: "http://api:8080"},
		{name: "falls back when unset", value: "", fallback: "http://localhost:8080", want: "http://localhost:8080"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("TEST_KEY", tt.value)

			if got := OrDefault("TEST_KEY", tt.fallback); got != tt.want {
				t.Errorf("OrDefault() = %q, want %q", got, tt.want)
			}
		})
	}
}
