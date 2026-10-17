# pismo-transactions-routine-service

Este serviço tem como objetivo implementar "Transactions Routine" — cuidando do ciclo de vida de contas e de transações, sendo extensível ao aplicar princípios de **Clean Architecture** e **arquitetura hexagonal**.

## Estrutura do Projeto

```
pismo-transactions-routine-service/
├── .gitignore
├── README.md
└── app/                       # Código da aplicação
    ├── cmd/                   # Pontos de entrada
    │   └── main.go            # Servidor HTTP (Gin) + health check
    ├── go.mod                 # Definição do módulo Go
    ├── go.sum                 # Checksums de dependências
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
| `make test` | Executa os testes e gera relatório de cobertura |

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
