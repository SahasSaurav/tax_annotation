.PHONY: build run clean test lint fmt vet

build:
	go build -o bin/taxrender ./cmd/taxrender/

run:
	go run ./cmd/taxrender/

clean:
	rm -rf bin/

test:
	go test ./...

lint:
	golangci-lint run ./...

fmt:
	gofmt -w .

vet:
	go vet ./...
