.PHONY: run tidy fmt test test-e2e compose-up compose-down up down remove clean migrate migrations

BASE_URL ?= http://localhost:8080


compose-up:
	@# Detect podman or fallback to docker
	@cmd=$$(command -v podman >/dev/null 2>&1 && echo podman || echo docker); \
	$$cmd compose up -d

up: compose-up

observability:
	@# Detect podman or fallback to docker
	@cmd=$$(command -v podman >/dev/null 2>&1 && echo podman || echo docker); \
	$$cmd compose up -d otel-collector prometheus jaeger grafana

compose-down:
	@# Detect podman or fallback to docker
	@cmd=$$(command -v podman >/dev/null 2>&1 && echo podman || echo docker); \
	$$cmd compose down

down: compose-down

remove:
	@# Detect podman or fallback to docker
	@cmd=$$(command -v podman >/dev/null 2>&1 && echo podman || echo docker); \
	$$cmd compose down --volumes --remove-orphans

clean: remove

run: up

run-local: observability db
	cd app && go run ./cmd

tidy:
	cd app && go mod tidy

fmt:
	cd app && go fmt ./...

test:
	cd app && go test ./... -coverprofile=coverage.out
	cd app && go tool cover -html=coverage.out

test-e2e: up
	@trap 'status=$$?; trap - EXIT; $(MAKE) -C "$(CURDIR)" down; down_status=$$?; if [ $$status -eq 0 ]; then exit $$down_status; else exit $$status; fi' EXIT; \
	cd app && BASE_URL="$(BASE_URL)" go test -tags=e2e ./tests/integration/... -v

migrations:
	@# Detect podman or fallback to docker
	@cmd=$$(command -v podman >/dev/null 2>&1 && echo podman || echo docker); \
	$$cmd compose run --rm migrations

migrate: migrations

db:
	@# Detect podman or fallback to docker
	@cmd=$$(command -v podman >/dev/null 2>&1 && echo podman || echo docker); \
	$$cmd compose up -d postgres cloudbeaver

