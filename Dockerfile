# Build stage
FROM golang:1.24-bookworm AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o bot ./src/cmd

# Runtime stage
FROM debian:bookworm-slim

# Install ca-certificates for HTTPS requests and sqlite3 for database
RUN apt-get update && apt-get install -y \
  ca-certificates \
  sqlite3 \
  && rm -rf /var/lib/apt/lists/*

# Copy the binary from builder stage
COPY --from=builder /app/bot /usr/local/bin/bot

# Copy the assets directory from builder stage
COPY --from=builder /app/src/assets /app/src/assets

# Set working directory to match the expected path structure
WORKDIR /app

# Run the binary
CMD ["bot"]
