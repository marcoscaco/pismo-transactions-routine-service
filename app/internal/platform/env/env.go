// Package env provides helpers to read configuration from the process environment.
// It has no knowledge of the domain, adapters or infrastructure layers.
package env

import "os"

// OrDefault returns the value of key, or fallback when it is unset or empty.
func OrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
