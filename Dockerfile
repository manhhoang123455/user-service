# Build stage
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
COPY .env .env
RUN go mod tidy
RUN go build -o user-service ./cmd/user-service

# Run stage
FROM alpine:3.15
WORKDIR /app
COPY --from=builder /app/user-service /app/user-service
COPY --from=builder /app/config /app/config
COPY --from=builder /app/.env /app/.env

CMD ["./user-service"]