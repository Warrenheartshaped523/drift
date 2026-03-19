.PHONY: build run install clean test lint setup help

BINARY := drift
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -s -w -X main.version=$(VERSION)

## build: compile the binary to ./drift
build:
	go build -ldflags "$(LDFLAGS)" -o $(BINARY) .

## run: build and run immediately
run: build
	./$(BINARY)

## install: install to GOPATH/bin
install:
	go install -ldflags "$(LDFLAGS)" .

## test: run all tests with race detector
test:
	go test -race -count=1 ./...

## lint: run golangci-lint
lint:
	golangci-lint run ./...

## clean: remove built artefacts
clean:
	rm -f $(BINARY)

## setup: download dependencies and tidy go.sum
setup:
	go mod download
	go mod tidy

## demo: record all demo GIFs with vhs (requires: brew install vhs)
demo:
	vhs demo/waveform.tape
	vhs demo/constellation.tape
	vhs demo/demo.tape

## help: print this help
help:
	@grep -E '^## ' Makefile | sed 's/## //' | column -t -s ':'
