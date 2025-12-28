#!/bin/sh
set -e

# Create storage directory if it doesn't exist
mkdir -p /app/storage

# If running as root, we can't change ownership (shouldn't happen with user directive)
# If running as a specific user (via docker-compose user directive), the directory
# will already have the correct permissions from the volume mount

# Execute the main application
exec "$@"
