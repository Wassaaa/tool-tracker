# Tool Tracker Makefile

all: build

stop: 
	docker compose down

################################################################################
# SETUP
################################################################################

deps:
	go mod download
	go mod tidy

deps-update:
	go get -u ./...
	go mod tidy

.PHONY: deps deps-update

################################################################################
# DEVELOPMENT
################################################################################

dev:
	docker-compose up -d && air

dev-build:
	docker-compose build

dev-logs:
	docker-compose logs -f

dev-stop:
	docker-compose down

dev-restart: dev-stop dev

.PHONY: dev dev-build dev-logs dev-stop dev-restart

################################################################################
# BUILD & TESTING
################################################################################

generate:
	go generate ./...

mocks:
	go generate ./cmd/api/internal/service

build: generate
	go build -o bin/tool-tracker cmd/api/main.go

build-clean:
	rm -rf bin/

test: generate
	go test ./...

test-coverage: generate
	go test -cover ./...

test-verbose: generate
	go test -v ./...

test-race: generate
	go test -race ./...

check-lint:
	golangci-lint run

check-vet:
	go vet ./...

check-fmt:
	go fmt ./...

check: check-fmt check-vet check-lint test

.PHONY: generate mocks build build-clean test test-coverage test-verbose test-race check-lint check-vet check-fmt check

################################################################################
# CLEAN
################################################################################

clean-mocks:
	rm -f cmd/api/internal/service/mocks/*.go

clean-artifacts:
	rm -rf bin/
	rm -f coverage.out
	rm -f tmp/main

clean: clean-artifacts
	docker-compose down --remove-orphans

clean-all: clean clean-mocks
	go clean -cache
	go clean -testcache

.PHONY: clean-mocks clean-artifacts clean clean-all

################################################################################
# HELP
################################################################################

help:
	@echo "Tool Tracker Development Commands"
	@echo ""
	@echo "Setup:"
	@echo "  deps            Download and tidy Go dependencies"
	@echo "  deps-update     Update all dependencies to latest versions"
	@echo ""
	@echo "Development:"
	@echo "  dev             Start Docker development environment with hot reload"
	@echo "  dev-build       Build Docker development containers"
	@echo "  dev-logs        Show development logs"
	@echo "  dev-stop        Stop development environment"
	@echo "  dev-restart     Restart development environment"
	@echo ""
	@echo "Build & Testing:"
	@echo "  generate        Generate all code (mocks, etc.)"
	@echo "  mocks           Generate only service mocks"
	@echo "  build           Build the application binary"
	@echo "  build-clean     Clean build artifacts"
	@echo "  test            Run all tests"
	@echo "  test-coverage   Run tests with coverage report"
	@echo "  test-verbose    Run tests with verbose output"
	@echo "  test-race       Run tests with race detection"
	@echo "  check-lint      Run Go linter"
	@echo "  check-vet       Run Go vet"
	@echo "  check-fmt       Format Go code"
	@echo "  check           Run all checks and tests"
	@echo ""
	@echo "Clean:"
	@echo "  clean-mocks     Remove generated mock files"
	@echo "  clean-artifacts Remove build artifacts"
	@echo "  clean           Stop services and clean artifacts"
	@echo "  clean-all       Deep clean including Go caches"

.PHONY: help