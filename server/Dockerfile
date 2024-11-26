FROM golang:1.23 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o server .

FROM debian:bookworm-slim

WORKDIR /app

# Install necessary runtime dependencies
RUN apt-get update && apt-get install -y \
    libc6 \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/server /app/server

EXPOSE 8080

CMD ["/app/server"]