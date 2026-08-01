.PHONY: test test-integration coverage coverage-check build run bench load-test

test:
	go test ./...

test-integration:
	go test -tags=integration ./internal/infrastructure/postgres/...

coverage:
	go test ./... -coverprofile=coverage.out -covermode=atomic
	go tool cover -func=coverage.out

coverage-check:
	go test ./... -coverprofile=coverage.out -covermode=atomic
	@go tool cover -func=coverage.out | findstr /C:"total:"

bench:
	go test -bench=. -benchmem ./internal/usecase/task/...

load-test:
	bash scripts/load_test.sh

build:
	go build -o bin/api ./cmd/api

run:
	go run ./cmd/api
