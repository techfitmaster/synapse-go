.PHONY: test test-all test-verbose lint coverage clean

test:
	go test ./... -race -cover -short

test-all:
	go test ./... -race -cover

test-verbose:
	go test ./... -race -cover -v -short

lint:
	golangci-lint run ./...

coverage:
	go test ./... -coverprofile=coverage.out -short
	go tool cover -html=coverage.out -o coverage.html

clean:
	rm -f coverage.out coverage.html
