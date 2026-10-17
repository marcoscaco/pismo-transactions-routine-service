# Stage 1: Build da aplicação em Go
FROM golang:1.25-alpine AS builder

WORKDIR /src

# Copia e baixa dependências
COPY app/go.mod app/go.sum ./
RUN go mod download

# Copia o código e gera o binário estático
COPY app/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /bin/api ./cmd/main.go

# Stage 2: Runtime minimalista com Distroless
FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

# Copia o binário estático do stage de build
COPY --from=builder /bin/api /app/api

# Usa o usuário não-root nativo do distroless por segurança
USER nonroot:nonroot

EXPOSE 8080

ENTRYPOINT ["/app/api"]