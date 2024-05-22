DOCKER ?= docker
TARGET ?= lordbasex/webssh
VERSION ?= latest

all:
	@echo "Building ${TARGET}"
	$(DOCKER) build -t $(TARGET):$(VERSION) -f Dockerfile .