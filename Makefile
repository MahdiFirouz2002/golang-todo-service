.PHONY: test test-integration coverage coverage-check build run

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

build:
	go build -o bin/api ./cmd/api

run:
	go run ./cmd/api
