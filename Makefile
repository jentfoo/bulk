export GO111MODULE = on

.PHONY: default test test-bulk test-store test-cover bench bench-bulk bench-store lint lint-bulk lint-store


test: test-bulk test-store

test-bulk:
	go test -race -cover ./...

test-store:
	cd store && go test -race -cover ./...

test-cover:
	go test -race -coverprofile=test.out ./... && go tool cover --html=test.out

bench: bench-bulk bench-store

bench-bulk:
	go test --benchmem -benchtime=10s -bench='Benchmark.*' -run='^$$'

bench-store:
	cd store && go test --benchmem -benchtime=10s -bench='Benchmark.*' -run='^$$'

lint: lint-bulk lint-store

lint-bulk:
	golangci-lint run --timeout=600s && go vet ./...

lint-store:
	cd store && golangci-lint run --timeout=600s && go vet ./...
