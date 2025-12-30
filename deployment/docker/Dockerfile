# Stage 1: Build Modules and Generate Code
FROM golang:1.25-bookworm AS builder

# Install system dependencies (libvips for image processing, sqlc for code gen)
RUN apt-get update && apt-get install -y --no-install-recommends \
    libvips-dev \
    && rm -rf /var/lib/apt/lists/*

# Install SQLC
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code (Optimized to avoid cache invalidation)
COPY cmd ./cmd
COPY internal ./internal
COPY sqlc.yaml schema.sql ./

# Generate SQLC code
RUN sqlc generate

# Build the application
# Optimized binary with stripped debug info
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-w -s" -o /app/bin/server ./cmd/api/main.go

# Stage 2: Build Frontend
FROM node:20-alpine AS frontend-builder
WORKDIR /app
COPY web/package*.json ./
# Use ci for deterministic builds
RUN npm ci
COPY web/ .
RUN npm run build

# Stage 3: Final Image
FROM debian:bookworm-slim

# Create non-root user
RUN groupadd -g 1001 appuser && \
    useradd -u 1001 -g appuser -s /bin/bash -m appuser

# Install runtime dependencies
RUN apt-get update && apt-get upgrade -y && apt-get install -y --no-install-recommends \
    libvips \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Create storage directory with correct permissions
RUN mkdir -p /app/storage && \
    chown -R appuser:appuser /app/storage

# Copy binary from builder with correct ownership
COPY --from=builder --chown=appuser:appuser /app/bin/server /app/server

# Copy frontend assets with correct ownership
COPY --from=frontend-builder --chown=appuser:appuser /app/dist /app/web/dist

# Copy entrypoint script with correct ownership and permissions
COPY --chown=appuser:appuser deployment/docker/entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh

# Expose port
EXPOSE 8080

# Environment variables
ENV PORT=8080
ENV APP_ENV=production

# Switch to non-root user
USER appuser

ENTRYPOINT ["/app/entrypoint.sh"]
CMD ["/app/server"]
