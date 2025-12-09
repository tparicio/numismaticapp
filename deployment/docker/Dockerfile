# Stage 1: Build Modules and Generate Code
FROM golang:1.25-bookworm AS builder

# Install system dependencies (libvips for image processing, sqlc for code gen)
RUN apt-get update && apt-get install -y \
    libvips-dev \
    && rm -rf /var/lib/apt/lists/*

# Install SQLC
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Generate SQLC code
RUN sqlc generate

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -o /app/bin/server ./cmd/api/main.go

# Stage 2: Build Frontend
FROM node:20-alpine AS frontend-builder
WORKDIR /app
COPY web/package*.json ./
# If package.json doesn't exist yet, this might fail. 
# We assume the web directory will be populated.
# For now, we'll use a conditional copy or assume the user follows the order.
# We will create a dummy package.json in the next steps if needed, 
# but for now let's assume standard flow.
RUN npm install
COPY web/ .
RUN npm run build

# Stage 3: Final Image
FROM debian:bookworm-slim

# Install runtime dependencies
RUN apt-get update && apt-get install -y \
    libvips \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/bin/server /app/server

# Copy frontend assets (served by Go Fiber or Nginx - assuming Fiber for simplicity here)
COPY --from=frontend-builder /app/dist /app/web/dist

# Copy migrations
# Copy migrations (Embedded in binary now, but keeping folder structure if needed for other tools)
# COPY --from=builder /app/migrations /app/migrations

# Expose port
EXPOSE 8080

# Environment variables
ENV PORT=8080
ENV APP_ENV=production

CMD ["/app/server"]
