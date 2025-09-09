FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o qrcode-generator-api ./cmd/api

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/qrcode-generator-api .
COPY --from=builder /app/.env .

EXPOSE 8080

CMD ["./qrcode-generator-api"]