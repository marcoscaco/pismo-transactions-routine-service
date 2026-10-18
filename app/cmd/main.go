package main

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"

	"github.com/marcoscaco/pismo-transactions-routine-service/internal/adapter/httpapi"
	"github.com/marcoscaco/pismo-transactions-routine-service/internal/infra/telemetry"
	"github.com/marcoscaco/pismo-transactions-routine-service/internal/platform/env"
	"github.com/marcoscaco/pismo-transactions-routine-service/internal/usecase"
)

const (
	defaultPort = "8080"
	serviceName = "pismo-transactions-routine-service"
)

func main() {
	ctx := context.Background()

	shutdown, err := telemetry.InitTelemetry(ctx, serviceName)
	if err != nil {
		log.Printf("Aviso: falha ao iniciar telemetria: %v", err)
	} else {
		defer shutdown()
	}

	metrics, err := httpapi.NewMetrics(otel.Meter(serviceName))
	if err != nil {
		log.Fatalf("falha ao criar instrumentos de métrica: %v", err)
	}

	// setup/wiring dos UseCases e Handlers
	accountUseCase := usecase.NewAccountUseCase()
	accountHandler := httpapi.NewAccountHandler(accountUseCase)

	// setup basico do router HTTP
	router := httpapi.NewRouter(
		httpapi.Handlers{
			Health:  httpapi.NewHealthHandler(),
			Account: accountHandler,
		},
		metrics,
	)

	if err := router.Run(serverAddress()); err != nil {
		log.Fatalf("falha ao iniciar servidor HTTP: %v", err)
	}
}

// serverAddress resolves the listen address from PORT, defaulting to 8080.
func serverAddress() string {
	return ":" + env.OrDefault("PORT", defaultPort)
}
