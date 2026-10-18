package main

import "testing"

func TestServerAddress(t *testing.T) {
	tests := []struct {
		name string
		port string
		want string
	}{
		{name: "default when unset", port: "", want: ":" + defaultPort},
		{name: "reads PORT", port: "9090", want: ":9090"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("PORT", tt.port)
			if got := serverAddress(); got != tt.want {
				t.Errorf("serverAddress() = %q, want %q", got, tt.want)
			}
		})
	}
}
