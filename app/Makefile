.PHONY: run tidy fmt

run:
	go run ./cmd

tidy:
	go mod tidy

fmt:
	go fmt ./...
test:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out


