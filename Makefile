
default: build

.PHONY: lint
lint:
	@echo "Linting..."
	@golangci-lint run

build:
	@echo "Building..."
	@goreleaser build --snapshot --clean --single-target

build/all:
	@echo "Building..."
	@goreleaser build --snapshot --clean

release/local:
	@echo "Releasing..."
	@goreleaser release --skip=validate --clean

clean:
	@echo "Cleaning..."
	@rm -rf dist
	@go mod tidy