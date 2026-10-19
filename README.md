# Pismo Transactions Routine Service

API HTTP em Go para criar e consultar contas e registrar transações, com persistência PostgreSQL e código organizado em camadas.

## Início rápido

Requer Docker com o plugin Docker Compose. Na raiz do repositório, inicie a stack:

```bash
docker compose up --build
```

Aguarde a API ficar disponível em `http://localhost:8080`. Em outro terminal, verifique a resposta com `curl http://localhost:8080/health`. Para encerrar a stack, use `docker compose down`. Para executar a aplicação no host ou consultar outros requisitos, veja [Desenvolvimento local](docs/DEVELOPMENT.md).

## Documentação de apoio

- [Contrato da API](docs/API.md): endpoints, formatos de requisição e resposta, erros e tipos iniciais de operação.
- [Desenvolvimento local](docs/DEVELOPMENT.md): requisitos, execução no host, configuração, testes, portas e observabilidade.
- [Arquitetura](docs/ARCHITECTURE.md): responsabilidades dos pacotes e fluxo de uma requisição.
