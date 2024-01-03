build:
	@echo "Building for Linux..."
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" main.go
	zip main.zip main
	@echo "Build for Linux complete."

all: build