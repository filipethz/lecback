# Etapa 1: Construir o aplicativo Go
FROM golang:1.21-alpine as builder

# Instalar dependências necessárias (como git) para o Go buscar pacotes
RUN apk add --no-cache git

WORKDIR /app
COPY . .

# Baixar dependências do Go
RUN go mod tidy
RUN go build -o app .

# Etapa 2: Executar o aplicativo Go
FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/app .

# Certifique-se de que a porta está sendo exposta corretamente
EXPOSE 8080

# Executar o binário gerado
CMD ["./app"]
