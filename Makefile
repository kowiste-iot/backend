# Makefile
BINARY_NAME=asset-service
VERSION?=1.0.0
BUILD_DIR=build
DOCKER_REPO=your-registry
SERVICE_NAME=asset-service

.PHONY: all build clean test coverage lint docker-build docker-push proto swagger swagger-validate help

all: clean lint test swagger build

## Build:
build: ## Build the binary
	@echo "Building..."
	@CGO_ENABLED=0 go build -ldflags="-X 'main.Version=${VERSION}'" -o ${BUILD_DIR}/${BINARY_NAME} ./cmd/api

build-linux: ## Build for Linux
	@echo "Building for Linux..."
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-X 'main.Version=${VERSION}'" -o ${BUILD_DIR}/${BINARY_NAME}-linux-amd64 ./cmd/api

up: ## 
	@echo "Create container..."
	@docker compose up -d

## Test:
test: ## Run unit tests
	@echo "Running tests..."
	@go test -v -race ./...

coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test -race -coverprofile=coverage.out -covermode=atomic ./...
	@go tool cover -html=coverage.out -o coverage.html

## Lint:
lint: ## Run linters
	@echo "Running linters..."
	@golangci-lint run ./...

## Proto:
proto: ## Generate Proto files
	@echo "Generating proto files..."
	@protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/*.proto

## Docker:
docker-build: ## Build docker image
	@echo "Building docker image..."
	@docker build -t ${DOCKER_REPO}/${SERVICE_NAME}:${VERSION} .
	@docker tag ${DOCKER_REPO}/${SERVICE_NAME}:${VERSION} ${DOCKER_REPO}/${SERVICE_NAME}:latest

docker-push: ## Push docker image
	@echo "Pushing docker image..."
	@docker push ${DOCKER_REPO}/${SERVICE_NAME}:${VERSION}
	@docker push ${DOCKER_REPO}/${SERVICE_NAME}:latest

## Clean:
clean: ## Clean build directory
	@echo "Cleaning..."
	@rm -rf ${BUILD_DIR}
	@rm -f coverage.out coverage.html


## Development:
dev: ## Run service in development mode
	@echo "Starting development server..."
	@go run ./cmd/server/main.go

## Tools:
tools: ## Install development tools
	@echo "Installing tools..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@go install github.com/swaggo/swag/cmd/swag@latest

nats: ## Run nats client
	@echo "Starting nats client..."
	@go run ./test/nats/subscriber.go

## Swagger:
swagger: swagger-validate ## Generate Swagger documentation
	@echo "Generating Swagger documentation..."
	@swag init -g cmd/api/main.go -o docs/swagger
	@echo "Swagger UI will be available at http://localhost:8080/swagger/index.html"

swagger-validate: ## Validate Swagger documentation
	@echo "Validating Swagger documentation..."
	@swag fmt -g cmd/api/main.go --exclude docs/swagger
	@swag fmt --validate -g cmd/api/main.go --exclude docs/swagger

swagger-serve: ## Serve Swagger documentation
	@echo "Starting Swagger UI server..."
	@docker run -p 8082:8080 -e SWAGGER_JSON=/docs/swagger/swagger.json -v $(shell pwd)/docs/swagger:/docs/swagger swaggerapi/swagger-ui

## Help:
help: ## Show this help
	@echo "Usage:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help