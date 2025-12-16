# Build Stage
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o ergo-api cmd/ergo-api/main.go

# Run Stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/ergo-api .
COPY --from=builder /app/policies ./policies
EXPOSE 8080
CMD ["./ergo-api"]
