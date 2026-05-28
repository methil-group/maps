# Stage 1: Build
FROM golang:1.26-alpine AS builder

# Install build dependencies (optional, but good for some Go packages)
RUN apk add --no-cache git

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application as a static binary
RUN CGO_ENABLED=0 GOOS=linux go build -o server_bin .

# Stage 2: Runtime
FROM alpine:latest

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata wget

# Copy the binary from the builder stage
COPY --from=builder /app/server_bin .

# Copy scripts and config
COPY scripts/download_maps.sh ./download_maps.sh
COPY MAPS.txt ./MAPS.txt
RUN chmod +x ./download_maps.sh

# Create a startup script
RUN echo -e '#!/bin/sh\n./download_maps.sh\nexec ./server_bin' > /app/start.sh
RUN chmod +x /app/start.sh

# Create the data directory
RUN mkdir -p /app/data

# Expose the default port
EXPOSE 8080

# Run the startup script
CMD ["/app/start.sh"]
