# pismo-transactions-routine-service

Este serviço tem como objetivo implementar "Transactions Routine" — cuidando do ciclo de vida de contas e de transações, sendo extensível ao aplicar princípios de **Clean Architecture** e **arquitetura hexagonal**.

## Estrutura do Projeto

```
pismo-transactions-routine-service/
├── .github/agents/          # Agentes de IA (Go Expert, Go Documentation Specialist, Go Solution Auditor)
├── .github/skills/          # Skills de IA (ex.: acquire-codebase-knowledge)
├── .gitignore
├── Dockerfile               # Imagem multi-estágio (golang:1.26-alpine -> distroless)
├── docker-compose.yml       # Orquestração local: api, otel-collector, prometheus, jaeger
├── README.md
└── app/                     # Código da aplicação (módulo Go)
    ├── cmd/                 # Pontos de entrada
    │   └── main.go          # Wiring: telemetria, métricas, handlers, router e servidor HTTP
    ├── go.mod               # Definição do módulo Go
    ├── go.sum               # Checksums de dependências
    ├── Makefile             # Comandos de conveniência
    ├── internal/
    │   ├── adapter/httpapi/ # Adapter HTTP (Gin)
    │   │   ├── health_handler.go   # GET /health -> {"status":"ok"}
    │   │   ├── health_handler_test.go
    │   │   ├── metrics.go          # Contador e histograma de requisições HTTP (OpenTelemetry)
    │   │   ├── metrics_test.go
    │   │   └── router.go           # Registro de rotas e middlewares
    │   ├── infra/telemetry/        # Infraestrutura de observabilidade (OpenTelemetry)
    │   │   ├── telemetry.go        # Inicialização do MeterProvider e TracerProvider (OTLP gRPC)
    │   │   └── telemetry_test.go
    │   └── platform/env/           # Leitura de variáveis de ambiente
    │       ├── env.go              # OrDefault(key, fallback)
    │       └── env_test.go
    └── tests/               # Testes
        └── integration/     # Testes de integração E2E (blackbox)
            └── healthcheck_e2e_test.go
```

## Como Executar

### Pré-requisitos

- Go 1.26+
- Make (para os comandos de conveniência)
- Docker e Docker Compose (para rodar com observabilidade via container)

### Iniciar o Serviço (local, sem Docker)

```bash
cd app
make run
```

O serviço inicia um servidor HTTP Gin na porta **8080** (configurável via variável de ambiente `PORT`) e expõe o endpoint `/health` para verificação de saúde.

### Iniciar com Docker Compose (observabilidade completa)

Sobe a API, o OpenTelemetry Collector, o Prometheus e o Jaeger:

```bash
docker compose up
```

Serviços expostos:

| Serviço | Porta | Descrição |
|---------|-------|-----------|
| api | 8080 | API Gin com telemetria |
| otel-collector | 4317 | Receptor OTLP gRPC |
| prometheus | 9090 | Coleta métricas do collector (8889) |
| jaeger | 16686 | Interface Web de rastreamento |

A API envia traces e métricas via OTLP para `otel-collector:4317`. O collector encaminha traces para o Jaeger e métricas para o Prometheus.

### Comandos Disponíveis (dentro de `app/`)

| Comando | Descrição |
|---------|-----------|
| `make run` | Inicia o serviço (`go run ./cmd`) |
| `make tidy` | Limpa dependências não utilizadas (`go mod tidy`) |
| `make fmt` | Formata o código com `gofmt` |
| `make test` | Executa os testes de unidade e gera relatório de cobertura |
| `make test-e2e` | Executa os testes de integração E2E (contra `BASE_URL`) |

### Verificando a Saúde

```bash
curl http://localhost:8080/health
```

Retorna `{"status":"ok"}`.

### Executando os Testes

```bash
cd app
make test
```

Isso gera um relatório de cobertura em `coverage.out` e exibe um HTML com `go tool cover -html=coverage.out`.

## Quick Start

1. Acesse a pasta do app:
   ```bash
   cd app
   ```

2. Inicie o serviço:
   ```bash
   make run
   ```

3. Verifique em `http://localhost:8080/health`

4. Execute os testes:
   ```bash
   make test
   ```

## Detalhes Técnicos

- **Linguagem**: Go 1.26.0
- **Framework**: Gin (HTTP web framework)
- **Observabilidade**: OpenTelemetry (traces e métricas via OTLP gRPC), Jaeger, Prometheus
- **Finalidade**: Gerenciar o ciclo de vida de contas e transações

## Configuração por Variáveis de Ambiente

Todas são lidas no pacote `app/internal/platform/env` (via `OrDefault`) ou diretamente via `os.Getenv`.

| Variável | Padrão | Descrição |
|----------|--------|-----------|
| `PORT` | `8080` | Porta na qual o servidor Gin escuta |
| `ENVIRONMENT` | `development` | Ambiente de execução (usado como atributo do recurso OTel) |
| `OTEL_EXPORTER_OTLP_ENDPOINT` | `otel-collector:4317` | Endpoint do coletor OTLP gRPC |
| `OTEL_EXPORTER_OTLP_INSECURE` | `true` | Se a conexão com o coletor deve ser insegura (sem TLS) |
| `OTEL_METRIC_EXPORT_INTERVAL` | `1s` | Intervalo de exportação de métricas (milissegundos; valor inválido usa o padrão) |
| `BASE_URL` | `http://localhost:8080` | URL base para os testes E2E blackbox |

## Contrato da API

### GET /health

Verificação de saúde do serviço.

**Resposta (200):**
```json
{"status":"ok"}
```

## Testes E2E Blackbox

Os testes de integração E2E estão em `app/tests/integration/` e seguem um modelo **blackbox**: não iniciam o próprio servidor, fazem requests HTTP reais contra um serviço externo configurado via `BASE_URL`. Isso permite reutilizar os mesmos testes contra qualquer ambiente (desenvolvimento, homologação, produção) sem modificações.

### Execução

```bash
# Rodar o serviço
cd app
make run

# Verificar healthcheck
curl http://localhost:8080/health

# Executar testes E2E
BASE_URL=http://localhost:8080 go test -run TestHealthCheckE2E -v
```

### Contra Outros Ambientes

```bash
# Contra homologação
BASE_URL=https://api-homologacao.empresa.com go test -run TestHealthCheckE2E -v

# Contra produção
BASE_URL=https://api.empresa.com go test -run TestHealthCheckE2E -v
```


## Contribuidores

- **Marcos Antonio Cordovil de Oliveira** (Caco): Principal contribuidor, arquitetura da solução e implementação.
- **Diversas ferramentas de IA generativa** (Vários modelos): Usadas principalmente na geração de testes.