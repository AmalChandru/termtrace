.PHONY: build test check fmt vet run clean

BINARY := bin/termtrace
CMD := ./cmd/termtrace

build:
	go build -o $(BINARY) $(CMD)

test:
	go test ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

check: fmt vet test

run: build
	./$(BINARY)

clean:
	rm -f $(BINARY)