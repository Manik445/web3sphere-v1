# Build stage
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Install build dependencies
RUN apk add --no-cache make gcc musl-dev

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o server ./cmd/server/main.go

# Production stage
FROM alpine:3.19

WORKDIR /app

# Add required runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Create non-root user
RUN adduser -D -g '' web3sphere
USER web3sphere

# Copy the binary from the builder stage
COPY --from=builder --chown=web3sphere:web3sphere /app/server /app/server

# Expose port
EXPOSE 8080

# Run the binary
CMD ["/app/server"]
