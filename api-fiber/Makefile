buildLinux:
	@echo "Building for Linux..."
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o apiserver .
	@echo "Build for Linux complete."

buildMac:
	@echo "Building for macOS..."
	env CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o apiserver .
	@echo "Build for macOS complete."

all: buildLinux buildMac