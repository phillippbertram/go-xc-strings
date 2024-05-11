

.PHONY: lint
lint:
	@echo "Linting..."
	@golangci-lint run

build:
	@echo "Building..."
	@goreleaser release --snapshot --clean

release/local:
	@echo "Releasing..."
	@goreleaser release --skip=validate --clean
