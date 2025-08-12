FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /app/online_subscriptions ./cmd/main.go

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/online_subscriptions .
COPY ./config ./config

ENV CONFIG_PATH=/app/config/local.yaml
EXPOSE 8080

ENTRYPOINT ["./online_subscriptions"]