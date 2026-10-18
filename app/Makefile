.PHONY: run tidy fmt test test-e2e

BASE_URL ?= http://localhost:8080

run:
	go run ./cmd

tidy:
	go mod tidy

fmt:
	go fmt ./...

test:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

test-e2e:
	BASE_URL=$(BASE_URL) go test -tags=e2e ./tests/integration/... -run TestHealthCheckE2E -v
