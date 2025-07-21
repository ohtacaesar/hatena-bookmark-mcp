.PHONY: build test clean run install lint lint-fix

build:
	go build -o bin/hatena-bookmark-mcp cmd/main.go

test:
	go test ./...

clean:
	rm -rf bin/

run:
	go run cmd/main.go

install:
	go install ./cmd/...

deps:
	go mod tidy
	go mod download

lint:
	golangci-lint run

lint-fix:
	golangci-lint run --fix