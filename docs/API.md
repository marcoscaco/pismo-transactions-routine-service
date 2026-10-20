# API HTTP

A API recebe e responde JSON. As rotas são registradas em [`router.go`](../app/internal/adapter/httpapi/router.go). Não há rota HTTP para consultar transações nem para listar tipos de operação.

## `GET /health`

Retorna `200 OK` quando o handler responde:

```bash
curl -i http://localhost:8080/health
```

```json
{"status":"ok"}
```

Esse endpoint não verifica a conexão com o PostgreSQL.

## `POST /accounts`

Cria uma conta. `document_number` é obrigatório e não pode ser vazio ou conter somente espaços em branco. O campo aceita até 14 caracteres, o suficiente para CPF (11 dígitos) ou CNPJ (14 dígitos). O esquema do banco aplica esse limite (`VARCHAR(14)`); o caso de uso não valida o tamanho antes da persistência.

**Requisição**

```json
{"document_number":"12345678900"}
```

```bash
curl -i -X POST http://localhost:8080/accounts \
  -H 'Content-Type: application/json' \
  -d '{"document_number":"12345678900"}'
```

**Resposta `201 Created`**

```json
{
  "account_id": 1,
  "document_number": "12345678900",
  "created_at": "<timestamp ISO 8601>"
}
```

## `GET /accounts/{id}`

Busca uma conta pelo ID. Em caso de sucesso (`200 OK`), retorna `account_id`, `document_number` e `created_at`, no mesmo formato da criação.

```bash
curl -i http://localhost:8080/accounts/1
```

## `POST /transactions`

Cria uma transação. `account_id` e `operation_type_id` devem ser IDs positivos; `amount` deve ser maior que zero. O multiplicador do tipo de operação determina o sinal gravado: `-1` para débito e `1` para crédito.

**Requisição — compra normal**

```json
{
  "account_id": 1,
  "operation_type_id": 1,
  "amount": "100"
}
```

```bash
curl -i -X POST http://localhost:8080/transactions \
  -H 'Content-Type: application/json' \
  -d '{"account_id":1,"operation_type_id":1,"amount":"100"}'
```

**Resposta `201 Created`**

```json
{
  "account_id": 1,
  "operation_type_id": 1,
  "amount": "-100"
}
```

**Requisição — crédito (valor numérico também é aceito)**

```bash
curl -i -X POST http://localhost:8080/transactions \
  -H 'Content-Type: application/json' \
  -d '{"account_id":1,"operation_type_id":4,"amount":60.0}'
```

Resposta: `{"account_id":1,"operation_type_id":4,"amount":"60"}`.

`amount` é representado como string decimal. A resposta não inclui o ID gerado nem a data do evento.

### Tipos de operação iniciais

O esquema inicial cadastra estes tipos:

| ID no banco vazio | Descrição | Multiplicador | Efeito |
|---:|---|---:|---|
| 1 | `Normal Purchase` | `-1` | Débito |
| 2 | `Purchase with installments` | `-1` | Débito |
| 3 | `Withdrawal` | `-1` | Débito |
| 4 | `Credit Voucher` | `1` | Crédito |

Os IDs são gerados pela coluna identity. A tabela mostra os valores gerados na inicialização de um banco vazio a partir do cadastro inicial; confirme os IDs disponíveis no banco antes de configurar clientes. Para consultar o catálogo armazenado:

```sql
SELECT id, description, amount_multiplier
FROM operation_types
ORDER BY id;
```

## Erros

Os erros retornados pelos handlers têm o formato `{"error":"mensagem"}`.

| Status | Quando ocorre |
|---|---|
| `400 Bad Request` | JSON inválido, documento ausente/em branco, IDs inválidos, valor não positivo ou tipo de operação inexistente. |
| `404 Not Found` | Conta inexistente consultada por ID ou referenciada por uma transação. |
| `500 Internal Server Error` | Falha interna não mapeada; a mensagem é genérica. |

O contrato é implementado pelos [handlers](../app/internal/adapter/httpapi) e pelo [mapeador de erros](../app/internal/adapter/httpapi/error_response.go).
