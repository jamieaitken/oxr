.DEFAULT_GOAL := ci

ci: lint test

test:
	@go test -race -failfast -covermode=atomic ./...

lint:
	@golangci-lint run ./...