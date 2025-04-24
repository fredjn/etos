.PHONY: all build clean test

# Build settings
BINARY_NAME=etosctl
PLUGIN_DIR=plugins
EXAMPLE_PLUGIN=example.so

all: build

build:
	go build -o $(BINARY_NAME)

build-plugin:
	cd $(PLUGIN_DIR)/example && go build -buildmode=plugin -o ../$(EXAMPLE_PLUGIN)

clean:
	rm -f $(BINARY_NAME)
	rm -f $(PLUGIN_DIR)/*.so

test:
	go test ./...

# Development tools
tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

lint:
	golangci-lint run

# Help target
help:
	@echo "Available targets:"
	@echo "  all          - Build everything (default target)"
	@echo "  build        - Build the main CLI binary"
	@echo "  build-plugin - Build the example plugin"
	@echo "  clean        - Remove build artifacts"
	@echo "  test         - Run tests"
	@echo "  tools        - Install development tools"
	@echo "  lint         - Run linter" 