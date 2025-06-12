# Build stage
FROM golang:1.23-alpine AS builder

# Instalar dependências necessárias
RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

# Copiar arquivos de dependências primeiro (para cache)
COPY go.mod go.sum ./
RUN go mod download

# Copiar código fonte
COPY . .

# Build da aplicação
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main cmd/server/main.go

# Final stage
FROM alpine:latest

# Instalar dependências mínimas
RUN apk --no-cache add ca-certificates sqlite wget tzdata && \
    rm -rf /var/cache/apk/*

# Criar usuário não-root
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

WORKDIR /app

# Copiar binário e templates
COPY --from=builder --chown=appuser:appgroup /app/main ./
COPY --from=builder --chown=appuser:appgroup /app/web ./web

# Criar diretório de dados
RUN mkdir -p /data && chown appuser:appgroup /data

# Trocar para usuário não-root
USER appuser

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/ || exit 1

CMD ["./main"]
