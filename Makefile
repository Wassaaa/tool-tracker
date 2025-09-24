# Tool Tracker Monorepo Makefile

.PHONY: help setup dev build test clean

all: dev

################################################################################
# SETUP
################################################################################

setup:
	pnpm install

################################################################################
# DEVELOPMENT
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

# Docker commands
docker-up:
	cd packages/backend && docker-compose up -d

docker-down:
	cd packages/backend && docker-compose down

docker-logs:
	cd packages/backend && docker-compose logs -f

docker-build:
	cd packages/backend && docker-compose build

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
