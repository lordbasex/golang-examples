# Variables
GO_CMD=/usr/local/go/bin/go
BUILD_FLAGS=-ldflags="-s -w"
OUTPUT_AMD64=asterisk-core-show-channels-amd64
OUTPUT_ARM64=asterisk-core-show-channels-aarch64
MAIN_FILE=main.go

.PHONY: all clean

# Regla para compilar ambos binarios
all: amd64 arm64

# Compilación para amd64
amd64:
	env GOOS=linux GOARCH=amd64 $(GO_CMD) build $(BUILD_FLAGS) -o $(OUTPUT_AMD64) $(MAIN_FILE)
	@echo "Built $(OUTPUT_AMD64)"

# Compilación para arm64
arm64:
	env GOOS=linux GOARCH=arm64 $(GO_CMD) build $(BUILD_FLAGS) -o $(OUTPUT_ARM64) $(MAIN_FILE)
	@echo "Built $(OUTPUT_ARM64)"

# Limpiar binarios generados
clean:
	rm -f $(OUTPUT_AMD64) $(OUTPUT_ARM64)
	@echo "Cleaned build artifacts"
