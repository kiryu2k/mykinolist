.PHONY: build
build:
	@go build -o ./bin/mykinolist ./cmd/mykinolist/main.go

.PHONY: run
run: build
	@./bin/mykinolist

.PHONY: test
test:
	@go test -v ./...

.PHONY: .install-linter
.install-linter:
	@[ -f ./bin/golangci-lint ] || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./bin v1.53.3

.PHONY: lint
lint: .install-linter
	@./bin/golangci-lint run ./...

.PHONY: lint-fast
lint-fast: .install-linter
	@./bin/golangci-lint run ./... --fast

DB_PORT := 5432
DB_USER := kirrryu
DB_PASSWORD := qwerty
DB_NAME := mykinolist
migrate-up:
	@migrate -path ./migrations -database 'postgres://$(DB_USER):$(DB_PASSWORD)@localhost:$(DB_PORT)/$(DB_NAME)?sslmode=disable' up

migrate-down:
	@migrate -path ./migrations -database 'postgres://$(DB_USER):$(DB_PASSWORD)@localhost:$(DB_PORT)/$(DB_NAME)?sslmode=disable' down

.DEFAULT_GOAL := run