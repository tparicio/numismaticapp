# Stage 1: Build Modules and Generate Code
FROM golang:1.25-alpine AS builder

# Install system dependencies (libvips for image processing, sqlc for code gen)
# build-base needed for CGO
RUN apk add --no-cache \
    vips-dev \
    build-base \
    git

# Install SQLC
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY cmd ./cmd
COPY internal ./internal
COPY sqlc.yaml schema.sql ./

# Generate SQLC code
RUN sqlc generate

# Build the application
# Optimized binary with stripped debug info
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-w -s" -o /app/bin/server ./cmd/api/main.go

# Stage 2: Build Frontend
FROM node:22-alpine AS frontend-builder
WORKDIR /app
COPY web/package*.json ./
# Use ci for deterministic builds
RUN npm ci
COPY web/ .
RUN npm run build

# Stage 3: Final Image
FROM alpine:3.23.2

# Create non-root user
RUN adduser -D -u 1001 appuser

# Install runtime dependencies
# bash is required for entrypoint.sh script
RUN apk add --no-cache \
    vips \
    ca-certificates \
    bash

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
