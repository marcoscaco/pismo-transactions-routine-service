# Desenvolvimento local

Este guia reúne os requisitos, opções de execução e comandos de desenvolvimento. A inicialização da stack completa está no [README](../README.md).

## Requisitos

- Go 1.26.0 ou versão compatível com `app/go.mod` para executar a aplicação no host.
- Docker com o plugin Docker Compose para banco, migrações e serviços auxiliares.
- `make` é opcional; fornece atalhos definidos no [`Makefile` da raiz](../Makefile).

## Iniciar a stack completa com `run.sh`

Na raiz do repositório:

```bash
./run.sh
```

O script inicia todos os serviços do Compose em segundo plano, com build. Ele usa `make up` quando o `make` está instalado; caso contrário, usa `podman compose` ou `docker compose`, nessa ordem de preferência. Se nenhum deles estiver disponível, encerra com erro.

Para parar a stack, use `make down`.

## Executar a API no host

Use este modo para rodar o código Go direto na sua máquina, com o restante da infraestrutura no Compose. Na raiz do repositório, execute:

```bash
make run-local
```

O alvo sobe os serviços de observabilidade e o PostgreSQL, e então executa `go run ./cmd` dentro de `app/`. Como o processo roda no host, as variáveis de conexão precisam apontar para `localhost`. Se preferir definir manualmente, use:

```bash
DB_HOST=localhost \
OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4317 \
OTEL_EXPORTER_OTLP_INSECURE=true \
make run-local
```

Caso o esquema ainda não esteja aplicado no banco, rode as migrações antes com `make migrate`.

O serviço `api` do Compose configura `DB_HOST=postgres` para usar o DNS da rede Docker. A execução no host usa `localhost` para alcançar a porta publicada do PostgreSQL.

## Configuração

A aplicação lê as variáveis abaixo no ponto de entrada ou na configuração de telemetria. Variáveis lidas por `env.OrDefault` usam o padrão quando estão vazias.

| Variável | Padrão | Descrição |
|---|---|---|
| `PORT` | `8080` | Porta HTTP. |
| `DB_HOST` | `localhost` no processo Go; `postgres` no serviço Compose `api` | Host do PostgreSQL. |
| `DB_PORT` | `5432` | Porta do PostgreSQL. |
| `DB_USER` | `pismo` | Usuário local do banco. |
| `DB_PASSWORD` | Valor local definido no código e no Compose | Credencial de desenvolvimento; configure um segredo apropriado fora desse ambiente. |
| `DB_NAME` | `pismo_db` | Nome do banco. |
| `DB_SSLMODE` | `disable` | Modo SSL do driver. |
| `ENVIRONMENT` | `development` | Atributo de ambiente do recurso OpenTelemetry. |
| `OTEL_EXPORTER_OTLP_ENDPOINT` | `otel-collector:4317` | Endpoint OTLP gRPC; no processo Go executado no host, use `localhost:4317`. |
| `OTEL_EXPORTER_OTLP_INSECURE` | `true` para o endpoint padrão | Configura transporte sem TLS. Um endpoint `https://` seleciona TLS por padrão; um valor explícito sobrescreve essa escolha. |
| `OTEL_METRIC_EXPORT_INTERVAL` | `1000` | Intervalo de exportação de métricas em milissegundos; valor inválido ou não positivo usa 1 segundo. |
| `BASE_URL` | `http://localhost:8080` | URL base usada pelo alvo E2E. |

Os valores de banco no Compose são destinados ao desenvolvimento local, não à produção.

## Testes e comandos

Execute os alvos `make` a partir da raiz do repositório. O comando Go direto deve ser executado dentro de `app/`:

| Comando | Diretório | Ação |
|---|---|---|
| `go test ./...` | `app/` | Executa os testes dos pacotes. |
| `make test` | raiz | Executa a suíte com cobertura e gera/exibe o relatório HTML. |
| `make test-e2e` | raiz | Sobe a stack local, executa com a build tag `e2e` e chama `make down` ao final, mesmo se o teste falhar. |
| `make fmt` | raiz | Formata os pacotes Go. |
| `make tidy` | raiz | Atualiza as dependências. |
| `./run.sh` | raiz | Inicia a stack completa em segundo plano (equivale a `make run`). |
| `make run` | raiz | Equivale a `make up`: sobe a stack do Compose em segundo plano. |
| `make run-local` | raiz | Sobe a infraestrutura e executa a API no host com `go run`. |

O teste E2E faz uma requisição HTTP real a `GET /health`. O alvo encerra os containers e a rede ao final; o volume persistente do PostgreSQL é mantido. `BASE_URL` pode ser sobrescrito, mas o alvo ainda inicia e encerra a stack local.

## Serviços locais e observabilidade

O Compose publica estes endereços no host:

| Serviço | Porta | Uso |
|---|---:|---|
| API | `8080` | Tráfego HTTP. |
| PostgreSQL | `5432` | Banco de dados. |
| OpenTelemetry Collector | `4317` | Receptor OTLP gRPC. |
| OpenTelemetry Collector | `8889` | Endpoint de scrape do Prometheus. |
| Prometheus | `9090` | Consulta de métricas. |
| Jaeger | `16686` | Interface de traces. |
| Grafana | `3000` | Dashboards e fontes de dados provisionadas. |
| CloudBeaver | `8081` | Interface web para o banco. |

O middleware HTTP exporta o contador `http_server_requests_total` e o histograma `http_server_request_duration_seconds`, com método, rota e status como atributos. O collector encaminha traces ao Jaeger e expõe métricas para o Prometheus; o Grafana é provisionado com esses dois serviços como fontes de dados.

Referências: [`Makefile`](../Makefile), [`app/cmd/main.go`](../app/cmd/main.go), [`telemetry.go`](../app/internal/infra/telemetry/telemetry.go) e [`docker-compose.yml`](../docker-compose.yml).
