FROM golang:1.23 AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o sse-server .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/sse-server .
EXPOSE 8080
CMD ["./sse-server"]