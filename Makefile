build:
	@go build -o bin/mykinolist ./cmd/mykinolist/main.go

run: build
	@./bin/mykinolist

test:
	@go test -v ./...

DB_PORT := 5432
DB_USER := kirrryu
DB_PASSWORD := qwerty
DB_NAME := mykinolist
migrate:
	@migrate -path ./migrations -database 'postgres://$(DB_USER):$(DB_PASSWORD)@localhost:$(DB_PORT)/$(DB_NAME)?sslmode=disable' up

.DEFAULT_GOAL := run