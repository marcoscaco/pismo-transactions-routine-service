.PHONY: run tidy fmt

run:
	go run ./cmd

tidy:
	go mod tidy

fmt:
	gofmt -w ./cmd
test:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out


