FROM golang:1.23-bookworm AS builder

WORKDIR /app

# Install required dependencies
RUN apt-get update && \
    apt-get install -y --no-install-recommends git ca-certificates && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Copy go.mod and go.sum first to leverage Docker caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the codebase
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o daisi-rest-postgres ./cmd/main.go

# Use a small debian slim image
FROM debian:bookworm-slim

WORKDIR /app

# Install CA certificates for HTTPS requests
RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates tzdata wget curl && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Copy binary from builder
COPY --from=builder /app/daisi-rest-postgres .

# Set environment variables
ENV GO_ENV=production

# Expose port
EXPOSE 80

# Run the service
CMD ["./daisi-rest-postgres"] 


