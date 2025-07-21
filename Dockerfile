FROM golang:1.24.5-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /rate-limiter ./cmd/server/

FROM alpine:latest

WORKDIR /app

COPY --from=builder /rate-limiter /app/rate-limiter
COPY .env.example /app/.env

EXPOSE 8080

CMD ["/app/rate-limiter"]