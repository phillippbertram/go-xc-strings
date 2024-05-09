
build:
	@echo "Building..."
	@go build -o bin/xcs

.PHONY: lint
lint:
	@echo "Linting..."
	@golangci-lint run