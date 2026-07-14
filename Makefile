.PHONY: build run test fmt vet tidy clean help

BINARY := persista-db
MAIN   := ./cmd/server
BIN    := bin/$(BINARY)

build:
	go build -o $(BIN) $(MAIN)

run: build
	./$(BIN)

test:
	go test ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

tidy:
	go mod tidy

clean:
	rm -rf bin/

help:
	@echo "Targets:"
	@echo "  build  - compile the server binary to $(BIN)"
	@echo "  run    - build and run the server"
	@echo "  test   - run all tests"
	@echo "  fmt    - format Go source files"
	@echo "  vet    - run go vet"
	@echo "  tidy   - tidy go.mod and go.sum"
	@echo "  clean  - remove build artifacts"
	@echo "  help   - show this help message"
