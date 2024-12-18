# Simple Makefile for a Go project
include .env
export $(shell sed 's/=.*//' .env)

# Build the application
all: build

build:
	@echo "Building..."


	@templ generate
	@go build -o api cmd/api/main.go
	@go build -o rpc cmd/rpc/main.go

# Run the application
run:
	@trap "kill 0" SIGINIT; \
	$(MAKE) run-rpc & \
	$(MAKE) run-api & \
	wait

run-api:
	@go run cmd/api/main.go

run-rpc:
	@go run cmd/rpc/main.go

# Create DB container
docker-run:
	@docker compose up -d

# Shutdown DB container
docker-down:
	@docker compose down

migrate:
	@dbmate migrate

migrate-create:
	@dbmate create

migrate-drop:
	@dbmate drop

migrate-up:
	@dbmate up

migrate-down:
	@dbmate down

generate:
	@sqlc generate

swag:
	@swag init -g internal/api/routes.go

protoc:
	@protoc --go_out=. --go-grpc_out=. proto/framer.proto

# Test the application
test:
	@echo "Testing..."
	@go test ./...


# Integrations Tests for the application
itest:
	@echo "Running integration tests..."
	@go test ./internal/pkg/database -v


# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

# Live Reload

watch:
	@if command -v air > /dev/null; then \
            air; \
            echo "Watching...";\
        else \
            read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
            if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
                go install github.com/air-verse/air@latest; \
                air; \
                echo "Watching...";\
            else \
                echo "You chose not to install air. Exiting..."; \
                exit 1; \
            fi; \
        fi


.PHONY: all build run test clean watch
