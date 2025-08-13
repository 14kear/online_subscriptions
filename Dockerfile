FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/online_subscriptions ./online_subscriptions/cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/bin/online_subscriptions .
COPY --from=builder /app/online_subscriptions/config ./config

ENV CONFIG_PATH=/app/config/local.yaml
EXPOSE 8080

CMD ["./online_subscriptions"]