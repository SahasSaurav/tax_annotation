.PHONY: build run clean test test-race test-verbose test-coverage lint fmt vet generate mocks mock-deps mock-gen

build:
	go build -o bin/taxrender ./cmd/taxrender/

run:
	go run ./cmd/taxrender/

clean:
	rm -rf bin/

test:
	go test ./...

test-race:
	go test -race ./...

test-verbose:
	go test -v ./...

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

lint:
	golangci-lint run ./...

fmt:
	gofmt -w .

vet:
	go vet ./...

mock-deps:
	go install go.uber.org/mock/mockgen@latest

mock-gen:
	mockgen -source=pkg/parser/interfaces.go -destination=mocks/mock_parser.go -package=mocks
	mockgen -source=pkg/formatter/interfaces.go -destination=mocks/mock_formatter.go -package=mocks
	mockgen -source=pkg/validator/interfaces.go -destination=mocks/mock_validator.go -package=mocks
	mockgen -source=pkg/render/interfaces.go -destination=mocks/mock_render.go -package=mocks
