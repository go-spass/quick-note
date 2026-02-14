BINARY := qn
BUILD_DIR := ./bin
CMD_DIR := ./cmd/qn

.PHONY: build test lint install clean

build:
	go build -o $(BUILD_DIR)/$(BINARY) $(CMD_DIR)

test:
	go test ./... -v

lint:
	golangci-lint run ./...

install:
	go install $(CMD_DIR)

clean:
	rm -rf $(BUILD_DIR)
