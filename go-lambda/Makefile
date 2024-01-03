build:
	@echo "Building for Linux..."
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o go-lambda main.go
	zip go-lambda.zip  go-lambda
	@echo "Build for Linux complete."

all: build