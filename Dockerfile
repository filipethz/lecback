# Etapa 1: Construir o aplicativo Go
FROM golang:1.20-alpine as builder

WORKDIR /app
COPY . .

RUN go mod tidy
RUN go build -o app .

# Etapa 2: Executar o aplicativo Go
FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/app .

EXPOSE 8080
CMD ["./app"]
