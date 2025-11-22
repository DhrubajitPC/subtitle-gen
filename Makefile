# Makefile for Video Subtitle Generator

# Variables
BINARY_NAME=video-subtitle-generator
GO_FILES=$(shell find . -name '*.go')

# Targets
.PHONY: all build test run clean

all: build

build:
	@echo "Building..."
	go build -o $(BINARY_NAME) main.go

test:
	@echo "Running tests..."
	go test -v ./...

run:
	@echo "Running application..."
	go run main.go

clean:
	@echo "Cleaning..."
	go clean
	rm -f $(BINARY_NAME)
