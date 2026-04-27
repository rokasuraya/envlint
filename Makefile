.PHONY: build test lint run clean

BINARY := envlint
CMD     := ./cmd/envlint

build:
	go build -o $(BINARY) $(CMD)

test:
	go test ./...

test-verbose:
	go test -v ./...

lint:
	golangci-lint run ./...

run: build
	./$(BINARY) -schema $(SCHEMA) -env $(ENV)

clean:
	rm -f $(BINARY)

# Example: make validate SCHEMA=internal/schema/testdata/example.schema.yaml ENV=.env
validate: build
	./$(BINARY) -schema $(SCHEMA) -env $(ENV)
