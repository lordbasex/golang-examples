build:
	@echo "Building for Linux..."
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" main.go
	zip main.zip main
	@echo "Build for Linux complete."

lambda:
	@echo "Create Function for Lambda"
	aws lambda create-function --function-name GO-LAMBDA-API-GATEWAY --zip-file fileb://main.zip --handler main --runtime go1.x --role arn:aws:iam:XXXXXXXXXXXXX
	@echo "Create Function for Lambda complete."

all: build