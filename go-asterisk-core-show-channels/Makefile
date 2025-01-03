# Variables
GO_CMD=/usr/local/go/bin/go
BUILD_FLAGS=-ldflags="-s -w"
MAIN_FILE=main.go

OUTPUT_LINUX_AMD64=go-asterisk-core-show-channels-linux-amd64
OUTPUT_LINUX_ARM64=go-asterisk-core-show-channels-linux-aarch64
OUTPUT_MACOS_AMD64=go-asterisk-core-show-channels-macos-amd64
OUTPUT_MACOS_ARM64=go-asterisk-core-show-channels-macos-arm64
OUTPUT_WINDOWS_AMD64=go-asterisk-core-show-channels-windows-amd64.exe
OUTPUT_WINDOWS_ARM64=go-asterisk-core-show-channels-windows-arm64.exe

.PHONY: all clean

# Regla para compilar todos los binarios
all: linux-amd64 linux-arm64 macos-amd64 macos-arm64 windows-amd64 windows-arm64

# Compilación para Linux amd64
linux-amd64:
	env GOOS=linux GOARCH=amd64 $(GO_CMD) build $(BUILD_FLAGS) -o $(OUTPUT_LINUX_AMD64) $(MAIN_FILE)
	@echo "Built $(OUTPUT_LINUX_AMD64)"

# Compilación para Linux arm64
linux-arm64:
	env GOOS=linux GOARCH=arm64 $(GO_CMD) build $(BUILD_FLAGS) -o $(OUTPUT_LINUX_ARM64) $(MAIN_FILE)
	@echo "Built $(OUTPUT_LINUX_ARM64)"

# Compilación para macOS amd64
macos-amd64:
	env GOOS=darwin GOARCH=amd64 $(GO_CMD) build $(BUILD_FLAGS) -o $(OUTPUT_MACOS_AMD64) $(MAIN_FILE)
	@echo "Built $(OUTPUT_MACOS_AMD64)"

# Compilación para macOS arm64
macos-arm64:
	env GOOS=darwin GOARCH=arm64 $(GO_CMD) build $(BUILD_FLAGS) -o $(OUTPUT_MACOS_ARM64) $(MAIN_FILE)
	@echo "Built $(OUTPUT_MACOS_ARM64)"

# Compilación para Windows amd64
windows-amd64:
	env GOOS=windows GOARCH=amd64 $(GO_CMD) build $(BUILD_FLAGS) -o $(OUTPUT_WINDOWS_AMD64) $(MAIN_FILE)
	@echo "Built $(OUTPUT_WINDOWS_AMD64)"

# Compilación para Windows arm64
windows-arm64:
	env GOOS=windows GOARCH=arm64 $(GO_CMD) build $(BUILD_FLAGS) -o $(OUTPUT_WINDOWS_ARM64) $(MAIN_FILE)
	@echo "Built $(OUTPUT_WINDOWS_ARM64)"

# Limpiar binarios generados
clean:
	rm -f $(OUTPUT_LINUX_AMD64) $(OUTPUT_LINUX_ARM64) \
	      $(OUTPUT_MACOS_AMD64) $(OUTPUT_MACOS_ARM64) \
	      $(OUTPUT_WINDOWS_AMD64) $(OUTPUT_WINDOWS_ARM64)
	@echo "Cleaned build artifacts"
