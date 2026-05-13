.PHONY: build test run clean lint

# Variáveis
BINARY_NAME=crom-pay
BUILD_DIR=bin
CMD_DIR=./cmd/crom-pay
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS=-ldflags "-s -w -X main.Version=$(VERSION)"

# Build local
build:
	@echo "⚡ Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)
	@echo "✅ Build completo: $(BUILD_DIR)/$(BINARY_NAME)"

# Executar
run:
	go run $(CMD_DIR)

# Testes
test:
	@echo "🧪 Rodando testes..."
	go test -v -race -count=1 ./...

# Testes com coverage
test-cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "📊 Coverage report: coverage.html"

# Lint
lint:
	@which golangci-lint > /dev/null 2>&1 || echo "⚠️  golangci-lint não instalado"
	golangci-lint run ./...

# Limpar
clean:
	@echo "🧹 Limpando..."
	rm -rf $(BUILD_DIR) coverage.out coverage.html
	@echo "✅ Limpo!"

# Cross-compilation
release:
	@echo "📦 Cross-compiling..."
	@mkdir -p dist
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-amd64 $(CMD_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-arm64 $(CMD_DIR)
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-amd64 $(CMD_DIR)
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-arm64 $(CMD_DIR)
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-windows-amd64.exe $(CMD_DIR)
	@echo "✅ Binários em dist/"
	@ls -lh dist/
