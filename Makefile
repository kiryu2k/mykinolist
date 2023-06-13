build:
	@go build -o bin/mykinolist ./cmd/mykinolist/main.go

run: build
	@./bin/mykinolist

test:
	@go test -v ./...

.DEFAULT_GOAL := run