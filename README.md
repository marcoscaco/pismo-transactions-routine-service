# Pismo Transactions Routine Service

API HTTP em Go para criar e consultar contas e registrar transações, com persistência PostgreSQL e código organizado em camadas.

## Início rápido

Requer Docker (ou Podman) com o plugin Compose. Na raiz do repositório, inicie a stack com o script:

```bash
./run.sh
```

O script usa `make up` quando o `make` está disponível e, caso contrário, `podman compose` ou `docker compose`, sempre em segundo plano. Aguarde a API ficar disponível em `http://localhost:8080` e verifique com `curl http://localhost:8080/health`.

Para encerrar a stack, use `make down`. Para executar a aplicação no host ou consultar outros requisitos, veja [Desenvolvimento local](docs/DEVELOPMENT.md).

## Documentação de apoio

- [Contrato da API](docs/API.md): endpoints, formatos de requisição e resposta, erros e tipos iniciais de operação.
- [Desenvolvimento local](docs/DEVELOPMENT.md): requisitos, execução no host, configuração, testes, portas e observabilidade.
- [Arquitetura](docs/ARCHITECTURE.md): responsabilidades dos pacotes e fluxo de uma requisição.
