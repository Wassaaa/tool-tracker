# Tool Tracker Monorepo Makefile

.PHONY: help setup dev build test clean

all: docker-up

################################################################################
# SETUP
################################################################################

setup:
	pnpm install

################################################################################
# DEVELOPMENT - LOCAL (legacy, may not work without CORS)
################################################################################

dev:
	pnpm dev

dev-frontend:
	pnpm dev:frontend

dev-backend:
	pnpm dev:backend

dev-stop:
	pnpm dev:stop

################################################################################
# DEVELOPMENT - CONTAINERIZED (recommended)
################################################################################

COMPOSE_FILE := docker-compose.yml

# Start containerized development environment
docker-up: ## Start containerized development environment with HTTPS
	@echo "🚀 Starting containerized development environment..."
	@if ! grep -q "tool-tracker.local" /etc/hosts; then \
		echo "� Note: tool-tracker.local not found in /etc/hosts"; \
		echo "   Using https://localhost as fallback"; \
	fi
	@docker-compose -f $(COMPOSE_FILE) up --build -d
	@echo ""
	@echo "✅ Services started!"
	@if grep -q "tool-tracker.local" /etc/hosts; then \
		echo "🌐 Main App: https://tool-tracker.local"; \
		echo "🔧 Backend API: https://tool-tracker.local/api"; \
		echo "📚 API Docs: https://tool-tracker.local/swagger/"; \
	else \
		echo "🌐 Main App: https://localhost"; \
		echo "🔧 Backend API: https://localhost/api"; \
		echo "📚 API Docs: https://localhost/swagger/"; \
	fi
	@echo "🗄️  Database Admin: http://localhost:9000"
	@echo ""
	@echo "⚠️  First time setup:"
	@echo "   1. Accept the self-signed certificate in your browser"
	@echo "   2. Or run 'make trust-ca' to install CA system-wide"

# Stop containerized development environment
docker-down: ## Stop containerized development environment
	@echo "🛑 Stopping development environment..."
	@docker-compose -f $(COMPOSE_FILE) down
	@echo "✅ Services stopped!"

# View containerized development logs
docker-logs: ## View containerized development logs (specify SERVICE=name for specific service)
	@docker-compose -f $(COMPOSE_FILE) logs -f $(SERVICE)

# Restart containerized services
docker-restart: ## Restart containerized services (specify SERVICE=name for specific service)
	@if [ -n "$(SERVICE)" ]; then \
		echo "🔄 Restarting $(SERVICE)..."; \
		docker-compose -f $(COMPOSE_FILE) restart $(SERVICE); \
	else \
		echo "🔄 Restarting all services..."; \
		docker-compose -f $(COMPOSE_FILE) restart; \
	fi

# Execute command in container
docker-exec: ## Execute command in container (specify SERVICE=name CMD="command")
	@docker-compose -f $(COMPOSE_FILE) exec $(or $(SERVICE),backend) $(CMD)

# Generate API docs and client in containers
docker-generate: ## Generate API docs and client in containers
	@echo "🔧 Generating API docs and client..."
	@docker-compose -f $(COMPOSE_FILE) exec backend make generate
	@docker-compose -f $(COMPOSE_FILE) exec frontend pnpm generate-api
	@echo "✅ Generation complete!"

# Check containerized environment status
docker-status: ## Check containerized environment status
	@echo "📊 Development Environment Status"
	@echo "=================================="
	@docker-compose -f $(COMPOSE_FILE) ps
	@echo ""
	@if grep -q "tool-tracker.local" /etc/hosts; then \
		echo "✅ tool-tracker.local configured"; \
	else \
		echo "💡 tool-tracker.local not found (using localhost)"; \
	fi

# Rebuild container images
docker-build: ## Rebuild all container images
	@echo "🔨 Building all services..."
	@docker-compose -f $(COMPOSE_FILE) build --no-cache
	@echo "✅ Build complete!"

# Clean up containerized environment
docker-clean: ## Clean up containerized environment completely
	@echo "🧹 Cleaning up development environment..."
	@docker-compose -f $(COMPOSE_FILE) down --remove-orphans
	@docker system prune -f
	@echo "✅ Cleanup complete!"

# Show containerized development help
docker-help: ## Show containerized development commands
	@echo "🛠️  Tool Tracker - Containerized Development Environment"
	@echo ""
	@echo "Quick Start:"
	@echo "  make docker-up     # Start everything with HTTPS"
	@echo "  make docker-down   # Stop everything"
	@echo "  make docker-clean  # Clean up completely"
	@echo ""
	@echo "Usage Examples:"
	@echo "  make docker-logs                              # View all logs"
	@echo "  make docker-logs SERVICE=backend             # View backend logs"
	@echo "  make docker-restart SERVICE=frontend         # Restart frontend"
	@echo "  make docker-exec SERVICE=frontend CMD='pnpm lint'  # Run command"
	@echo ""
	@echo "🌐 After 'make docker-up', visit: https://tool-tracker.local"

################################################################################
# DEV CERTS - CA Certificate Management
################################################################################
# Windows (Admin PowerShell) trust example:
#   certutil -addstore -f Root caddy-docker-root.crt
# Windows untrust example:
#   certutil -delstore Root "Caddy Local Authority - 2025 ECC Root"
################################################################################

CA_DOCKER_CERT := ./caddy-docker-root.crt
CA_DOCKER_ROOT := /data/caddy/pki/authorities/local/root.crt

# Extract CA certificate from Caddy container
get-ca: ## Extract Caddy's CA certificate for HTTPS trust
	@echo "🔐 Exporting Caddy CA root certificate..."
	@docker cp $$(docker-compose -f $(COMPOSE_FILE) ps -q caddy):$(CA_DOCKER_ROOT) $(CA_DOCKER_CERT) 2>/dev/null || \
		(echo "❌ Failed to extract CA cert. Is Caddy running? Try 'make docker-up' first." && exit 1)
	@echo "✅ Exported: $(CA_DOCKER_CERT)"
	@echo "📋 Next: 'make trust-ca' to install system-wide"

# Install CA certificate system-wide
trust-ca: get-ca ## Install Caddy's CA certificate system-wide (requires sudo)
	@echo "🔐 Installing CA certificate system-wide..."
	@OS=$$(uname); \
	if [ "$$OS" = "Darwin" ]; then \
		sudo security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain "$(CA_DOCKER_CERT)"; \
	elif [ "$$OS" = "Linux" ]; then \
		sudo cp "$(CA_DOCKER_CERT)" /usr/local/share/ca-certificates/caddy-docker-root.crt; \
		sudo update-ca-certificates; \
	else \
		echo "⚠️  Unsupported OS: $$OS"; \
		echo "Manual install: Import $(CA_DOCKER_CERT) as trusted root CA"; \
	fi
	@echo "✅ CA certificate installed! Restart your browser."
	@rm -f $(CA_DOCKER_CERT)

# Remove CA certificate from system
untrust-ca: ## Remove Caddy's CA certificate from system (requires sudo)
	@echo "🗑️  Removing Caddy CA certificate..."
	@OS=$$(uname); \
	if [ "$$OS" = "Darwin" ]; then \
		sudo security delete-certificate -c "Caddy Local Authority" /Library/Keychains/System.keychain 2>/dev/null || true; \
	elif [ "$$OS" = "Linux" ]; then \
		sudo rm -f /usr/local/share/ca-certificates/caddy-docker-root.crt; \
		sudo update-ca-certificates; \
	fi
	@echo "✅ CA certificate removed! Restart your browser."

################################################################################
# BUILD
################################################################################

build:
	pnpm build

build-frontend:
	pnpm build:frontend

build-backend:
	pnpm build:backend

################################################################################
# TESTING
################################################################################

test:
	pnpm test

test-frontend:
	pnpm test:frontend

test-backend:
	pnpm test:backend

################################################################################
# CLEAN
################################################################################

clean:
	pnpm clean

clean-go:
	cd packages/backend && rm -rf bin/ && rm -f coverage.out && rm -f coverage.html && rm -f tmp/main

clean-docker:
	cd packages/backend && docker-compose down --remove-orphans

clean-all: clean clean-go clean-docker
	rm -rf node_modules
	rm -rf packages/*/node_modules
	rm -rf packages/backend/tmp
	cd packages/backend && go clean -cache && go clean -testcache

################################################################################
# BACKEND SPECIFIC
################################################################################

# Go dependencies
deps:
	cd packages/backend && go mod download && go mod tidy

deps-update:
	cd packages/backend && go get -u ./... && go mod tidy

# Code generation
generate: mocks swagger api-client

mocks:
	cd packages/backend && go generate ./cmd/api/internal/service

swagger:
	cd packages/backend && swag init -g cmd/api/main.go -o docs

api-client: swagger
	cd packages/frontend && pnpm generate-api

# Backend build
build-go:
	cd packages/backend && go build -o bin/tool-tracker cmd/api/main.go

# Backend testing
test-go:
	cd packages/backend && go test ./...

test-repo:
	cd packages/backend && go test ./cmd/api/internal/repo -v

test-service:
	cd packages/backend && go test ./cmd/api/internal/service -v

test-domain:
	cd packages/backend && go test ./cmd/api/internal/domain -v

test-coverage:
	cd packages/backend && go test -cover ./...

test-coverage-html:
	cd packages/backend && go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: packages/backend/coverage.html"

test-verbose:
	cd packages/backend && go test -v ./...

test-race:
	cd packages/backend && go test -race ./...

test-integration:
	cd packages/backend && go test -tags=integration ./...

test-unit:
	cd packages/backend && go test -short ./...

test-benchmark:
	cd packages/backend && go test -bench=. ./...

test-watch:
	@command -v entr >/dev/null 2>&1 || { echo "entr is required for watch mode. Install it first."; exit 1; }
	cd packages/backend && find . -name "*.go" | entr -c go test ./...

# Code quality
check-lint:
	cd packages/backend && golangci-lint run

check-vet:
	cd packages/backend && go vet ./...

check-fmt:
	cd packages/backend && go fmt ./...

check-go: check-fmt check-vet check-lint test-go



################################################################################
# HELP
################################################################################

help:
	@echo "Tool Tracker Monorepo Commands"
	@echo ""
	@echo "Setup:"
	@echo "  setup           Install all dependencies"
	@echo "  deps            Download and tidy Go dependencies"
	@echo "  deps-update     Update all Go dependencies"
	@echo ""
	@echo "Development:"
	@echo "  dev             Start development environment (both frontend & backend)"
	@echo "  dev-frontend    Start only frontend"
	@echo "  dev-backend     Start only backend"
	@echo "  dev-stop        Stop development environment"
	@echo ""
	@echo "Build:"
	@echo "  build           Build all packages"
	@echo "  build-frontend  Build only frontend"
	@echo "  build-backend   Build only backend"
	@echo "  build-go        Build Go binary"
	@echo ""
	@echo "Testing:"
	@echo "  test            Run all tests"
	@echo "  test-frontend   Run frontend tests"
	@echo "  test-backend    Run backend tests"
	@echo "  test-go         Run Go tests"
	@echo "  test-coverage   Run Go tests with coverage"
	@echo "  test-verbose    Run Go tests with verbose output"
	@echo "  test-race       Run Go tests with race detection"
	@echo "  test-watch      Watch Go tests (requires entr)"
	@echo ""
	@echo "Code Generation:"
	@echo "  generate        Generate all code (mocks, swagger, frontend API)"
	@echo "  mocks           Generate service mocks"
	@echo "  swagger         Generate Swagger documentation"
	@echo "  api-client      Generate frontend API client"
	@echo ""
	@echo "Code Quality:"
	@echo "  check-go        Run all Go checks (fmt, vet, lint, test)"
	@echo "  check-fmt       Format Go code"
	@echo "  check-vet       Run Go vet"
	@echo "  check-lint      Run Go linter (requires golangci-lint)"
	@echo ""
	@echo "Docker:"
	@echo "  docker-up       Start Docker containers"
	@echo "  docker-down     Stop Docker containers"
	@echo "  docker-logs     Show Docker logs"
	@echo "  docker-build    Build Docker containers"
	@echo ""
	@echo "Clean:"
	@echo "  clean           Clean build artifacts"
	@echo "  clean-go        Clean Go build artifacts"
	@echo "  clean-docker    Stop Docker containers"
	@echo "  clean-all       Clean everything including caches"
