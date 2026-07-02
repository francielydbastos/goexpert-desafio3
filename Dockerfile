# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/ordersystem ./cmd/ordersystem

# Runtime stage
FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/ordersystem /app/ordersystem
COPY --from=builder /app/migrations /app/migrations

EXPOSE 8000 8080 50051

ENTRYPOINT ["/app/ordersystem"]
