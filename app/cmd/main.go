package main

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"

	"github.com/marcoscaco/pismo-transactions-routine-service/internal/adapter/httpapi"
	"github.com/marcoscaco/pismo-transactions-routine-service/internal/infra/database/postgres"
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

	// Configuração do banco de dados PostgreSQL
	// nunca usariamos valores hardcoded em produção, mas para fins de teste e desenvolvimento, vamos usar valores default
	sqlDB, err := postgres.NewPostgresDB(postgres.Config{
		Host:     env.OrDefault("DB_HOST", "localhost"),
		Port:     env.OrDefault("DB_PORT", "5432"),
		User:     env.OrDefault("DB_USER", "pismo"),
		Password: env.OrDefault("DB_PASSWORD", "pismo_pass"),
		DBName:   env.OrDefault("DB_NAME", "pismo_db"),
		SSLMode:  env.OrDefault("DB_SSLMODE", "disable"),
	})
	if err != nil {
		log.Fatalf("falha ao conectar ao banco de dados: %v", err)
	}
	defer sqlDB.Close()

	accountRepository := postgres.NewAccountRepository(sqlDB)
	accountUseCase := usecase.NewAccountUseCase(accountRepository)
	accountHandler := httpapi.NewAccountHandler(accountUseCase)

	transactionRepository := postgres.NewTransactionRepository(sqlDB)
	operationTypeRepository := postgres.NewOperationTypeRepository(sqlDB)
	transactionUseCase := usecase.NewTransactionUseCase(transactionRepository, operationTypeRepository)
	transactionHandler := httpapi.NewTransactionHandler(transactionUseCase)

	// setup basico do router HTTP
	router := httpapi.NewRouter(
		httpapi.Handlers{
			Health:      httpapi.NewHealthHandler(),
			Account:     accountHandler,
			Transaction: transactionHandler,
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
