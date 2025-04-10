# Build stage
FROM golang:1.23.3 AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o sse-server main.go

# Final stage
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/sse-server .
EXPOSE 8081
CMD ["./sse-server"]