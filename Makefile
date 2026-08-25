.PHONY: build run clean test lint fmt vet

build:
	go build -o bin/taxannotation main.go

run:
	go run main.go

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
