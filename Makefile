.DEFAULT_GOAL := help
.PHONY: help run stop logs lint test-unit test-unit-coverage generate clean tidy get docker-init docker-build docker-push deploy restart ci-local install-hooks vet build-web list-models

# Variables
DOCKER_COMPOSE = docker compose -f deployment/docker/docker-compose.yml
DOCKER_RUN_GO = docker run --rm -v $$(pwd):/app -w /app golang:1.25-bookworm
DOCKER_RUN_NODE = docker run --rm -v $$(pwd)/web:/app -w /app node:20-alpine
DOCKER_IMAGE ?= tparicio/numismaticapp
DOCKER_TAG ?= latest

# Colors
COLOR_RESET = \033[0m
COLOR_BOLD = \033[1m
COLOR_GREEN = \033[32m
COLOR_YELLOW = \033[33m
COLOR_BLUE = \033[34m
COLOR_CYAN = \033[36m

## ----------------------------------------------------------------------
## 🚀 NumismaticApp Makefile
## ----------------------------------------------------------------------

## 🐳 Docker Operations
run: ## Start the application in detached mode
	@echo "$(COLOR_BLUE)🚀 Starting application...$(COLOR_RESET)"
	$(DOCKER_COMPOSE) up --build -d
	@echo "$(COLOR_GREEN)✅ Application started!$(COLOR_RESET)"
	@echo "   Frontend: http://localhost:8080"
	@echo "   API:      http://localhost:8080/api/v1"

logs: ## View application logs
	@echo "$(COLOR_BLUE)📋 Tailing logs... (Ctrl+C to exit)$(COLOR_RESET)"
	$(DOCKER_COMPOSE) logs -f

stop: ## Stop the application
	@echo "$(COLOR_BLUE)🛑 Stopping application...$(COLOR_RESET)"
	$(DOCKER_COMPOSE) down
	@echo "$(COLOR_GREEN)✅ Application stopped.$(COLOR_RESET)"

restart: ## Restart the application (rebuilds app to apply migrations)
	@echo "$(COLOR_BLUE)🔄 Restarting application...$(COLOR_RESET)"
	$(DOCKER_COMPOSE) up -d --build app
	@echo "$(COLOR_GREEN)✅ Application restarted.$(COLOR_RESET)"
	@echo "$(COLOR_BLUE)📋 Tailing logs... (Ctrl+C to exit)$(COLOR_RESET)"
	$(DOCKER_COMPOSE) logs -f

docker-init: ## Initialize Docker Buildx for multi-arch builds
	@echo "$(COLOR_BLUE)🛠️  Initializing Docker Buildx...$(COLOR_RESET)"
	docker buildx create --use || true
	@echo "$(COLOR_GREEN)✅ Buildx initialized.$(COLOR_RESET)"

docker-build: ## Build Docker image locally (AMD64)
	@echo "$(COLOR_BLUE)🐳 Building Docker image $(DOCKER_IMAGE):$(DOCKER_TAG)...$(COLOR_RESET)"
	docker build -f deployment/docker/Dockerfile -t $(DOCKER_IMAGE):$(DOCKER_TAG) .
	@echo "$(COLOR_GREEN)✅ Image built: $(DOCKER_IMAGE):$(DOCKER_TAG)$(COLOR_RESET)"

docker-push: ## Build & Push Image (AMD64) to DockerHub
	@echo "$(COLOR_BLUE)🚀 Building and Pushing image $(DOCKER_IMAGE):$(DOCKER_TAG)...$(COLOR_RESET)"
	docker buildx build --platform linux/amd64 \
		-t $(DOCKER_IMAGE):$(DOCKER_TAG) \
		-f deployment/docker/Dockerfile \
		--push .
	@echo "$(COLOR_GREEN)✅ Image pushed!$(COLOR_RESET)"

deploy: docker-init docker-push ## Initialize buildx and push image

## 🛠️  Development
ci-local: lint vet test-unit ## Run local CI checks (Lint, Vet, Test)

install-hooks: ## Install git hooks
	@echo "$(COLOR_BLUE)🪝 Installing git hooks...$(COLOR_RESET)"
	cp deployment/scripts/git-hooks/commit-msg .git/hooks/commit-msg
	cp deployment/scripts/git-hooks/pre-push .git/hooks/pre-push
	chmod +x .git/hooks/commit-msg .git/hooks/pre-push
	@echo "$(COLOR_GREEN)✅ Git hooks installed!$(COLOR_RESET)"

vet: ## Run go vet
	@echo "$(COLOR_BLUE)🔍 Running go vet...$(COLOR_RESET)"
	go vet ./...

lint: ## Run linters (Go & Vue)
	@echo "$(COLOR_BLUE)🔍 Running linters...$(COLOR_RESET)"
	docker run --rm -v $$(pwd):/app -w /app golangci/golangci-lint:latest-alpine golangci-lint run -v
	cd web && npm run lint

build-web: ## Build frontend (in Docker)
	@echo "$(COLOR_BLUE)🏗️  Building frontend...$(COLOR_RESET)"
	$(DOCKER_RUN_NODE) npm install
	$(DOCKER_RUN_NODE) npm run build
	@echo "$(COLOR_GREEN)✅ Frontend built.$(COLOR_RESET)"

generate: ## Generate SQLC code
	@echo "$(COLOR_BLUE)⚙️  Generating SQLC code...$(COLOR_RESET)"
	docker run --rm -v $$(pwd):/src -w /src kjconroy/sqlc generate
	@echo "$(COLOR_GREEN)✅ Code generated.$(COLOR_RESET)"

list-models: ## List available Gemini models using curl
	@if [ -f .env.local ]; then \
		export $$(cat .env.local | xargs); \
		curl -s "https://generativelanguage.googleapis.com/v1beta/models?key=$$GEMINI_API_KEY" | grep '"name":' | sed 's/.*"name": "models\/\([^"]*\)".*/\1/'; \
	else \
		echo "Error: .env.local file not found"; \
	fi

# Test packages (exclude infrastructure, mocks, cmd, api from unit tests)
TEST_PKGS = $(shell go list ./... | grep -v 'internal/infrastructure' | grep -v 'internal/api' | grep -v 'cmd' | grep -v 'mocks')

## 🧪 Testing
test-unit: ## Run unit tests
	@echo "$(COLOR_BLUE)🧪 Running unit tests...$(COLOR_RESET)"
	go test -v $(TEST_PKGS)

test-unit-coverage: ## Run unit tests with coverage
	@echo "$(COLOR_BLUE)🧪 Running unit tests with coverage...$(COLOR_RESET)"
	@mkdir -p reports/tests
	go test -v -coverprofile=reports/tests/coverage.out $(TEST_PKGS)
	go tool cover -func=reports/tests/coverage.out

## 📦 Dependencies
tidy: ## Run go mod tidy (in Docker)
	@echo "$(COLOR_BLUE)🧹 Tidy up modules...$(COLOR_RESET)"
	$(DOCKER_RUN_GO) go mod tidy

get: ## Run go get (in Docker). Usage: make get PKG=...
	@echo "$(COLOR_BLUE)📥 Getting package $(PKG)...$(COLOR_RESET)"
	$(DOCKER_RUN_GO) go get $(PKG)

## 🧹 Cleanup
clean: ## Clean build artifacts
	@echo "$(COLOR_BLUE)🧹 Cleaning artifacts...$(COLOR_RESET)"
	rm -rf bin/
	rm -rf web/dist/
	@echo "$(COLOR_GREEN)✅ Cleaned.$(COLOR_RESET)"

## ❓ Help
help: ## Show this help message
	@echo ""
	@echo "$(COLOR_BOLD)🚀 NumismaticApp Manager$(COLOR_RESET)"
	@echo ""
	@echo "Usage: make $(COLOR_YELLOW)<target>$(COLOR_RESET)"
	@awk 'BEGIN {FS = ":.*##"; printf ""} \
		/^[a-zA-Z_-]+:.*?##/ { printf "  $(COLOR_YELLOW)%-20s$(COLOR_RESET) %s\n", $$1, $$2 } \
		/^## .*$$/ { printf "\n$(COLOR_CYAN)%s$(COLOR_RESET)\n", substr($$0, 4) } \
		' $(MAKEFILE_LIST)
	@echo ""
