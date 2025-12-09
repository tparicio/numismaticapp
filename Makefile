.PHONY: run stop build lint test generate clean deps

# Docker Compose
run:
	docker-compose -f deployment/docker/docker-compose.yml up --build

stop:
	docker-compose -f deployment/docker/docker-compose.yml down

# Local Development (Requires local tools)
lint:
	golangci-lint run
	cd web && npm run lint

test:
	go test -v ./...

generate:
	# Run sqlc using docker to avoid local dependency issues
	docker run --rm -v $$(pwd):/src -w /src kjconroy/sqlc generate

clean:
	rm -rf bin/
	rm -rf web/dist/

# Helper to run go mod tidy inside docker
tidy:
	docker run --rm -v $$(pwd):/app -w /app golang:1.25-bookworm go mod tidy

# Helper to run go get inside docker
# Usage: make get PKG=github.com/some/package
get:
	docker run --rm -v $$(pwd):/app -w /app golang:1.25-bookworm go get $(PKG)
