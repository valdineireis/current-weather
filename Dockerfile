# Estágio de Build
FROM golang:latest AS builder
RUN apt-get update && apt-get install -y ca-certificates
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
# Constrói o executável estaticamente, otimizado para Linux
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Estágio Final - Imagem de Produção
FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
WORKDIR /
COPY --from=builder /app/main .
EXPOSE 8080
ENTRYPOINT ["/main"]