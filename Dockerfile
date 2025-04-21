FROM golang:1.24-alpine AS builder

# Instala dependências mínimas
RUN apk add --no-cache git

WORKDIR /app

# Copia os arquivos go
COPY go.mod go.sum ./
RUN go mod download

# Copia o restante da aplicação
COPY . .

# Compila a aplicação (binário estático)
RUN go build -o orchestrator ./cmd/orchestrator/

# Etapa 2: Imagem final minimalista
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/orchestrator .

ENTRYPOINT ["./orchestrator"]
