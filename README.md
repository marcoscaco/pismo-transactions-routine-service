# pismo-transactions-routine-service

Este serviço tem como objetivo implementar "Transactions Routine" — cuidando do ciclo de vida de contas e de transações, sendo extensível ao aplicar princípios de **Clean Architecture** e **arquitetura hexagonal**.

## Estrutura do Projeto

```
pismo-transactions-routine-service/
├── .gitignore
├── README.md
└── app/                       # Código da aplicação
    ├── cmd/                   # Pontos de entrada
    │   └── main.go            # Wiring: handlers, router e servidor HTTP
    ├── internal/
    │   └── adapter/
    │       └── httpapi/       # Adapter HTTP (Gin)
    │           ├── health_handler.go
    │           ├── health_handler_test.go
    │           └── router.go
    ├── go.mod                 # Definição do módulo Go
    ├── go.sum                 # Checksums de dependências
    ├── tests/                 # Testes
    │   └── integration/       # Testes de integração E2E
    │       └── healthcheck_e2e_test.go
    └── Makefile               # Comandos de conveniência
```

## Como Executar

### Pré-requisitos

- Go 1.25+
- Make (para os comandos de conveniência)

### Iniciar o Serviço

```bash
cd app
make run
```

O serviço inicia um servidor HTTP Gin na porta **8080** e expõe o endpoint `/health` para verificação de saúde.

### Comandos Disponíveis

| Comando | Descrição |
|---------|-----------|
| `make run` | Inicia o serviço em modo de desenvolvimento |
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

- **Linguagem**: Go 1.25+
- **Framework**: Gin (HTTP web framework)
- **Finalidade**: Gerenciar o ciclo de vida de contas e transações

## Testes E2E Blackbox

Os testes de integração E2E são localizados em `app/tests/integration/` e seguem um modelo **blackbox**: eles não iniciam o próprio servidor, mas fazem requests HTTP reais contra um serviço externo configurado via variável de ambiente. Isso permite reutilizar os mesmos testes contra qualquer ambiente (desenvolvimento, homologação, produção) sem modificações.

### Configuração

O teste de healthcheck usa uma única variável de ambiente:

- **`BASE_URL`**: Endpoint base do serviço (padrão: `http://localhost:8080`)
  - Exemplos: `http://localhost:8080`, `http://api-homologacao.empresa.com`, `http://meu-servico:80`

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
