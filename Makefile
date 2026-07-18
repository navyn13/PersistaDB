.PHONY: build run dev peer test fmt vet tidy clean help

BINARY := persista-db
MAIN   := ./cmd/server
BIN    := bin/$(BINARY)
PEER_HOST := localhost
PEER_PORT := 9000

build:
	go build -o $(BIN) $(MAIN)

run: build
	./$(BIN)

# Hot-reload server with Air (rebuilds + restarts on .go save)
dev:
	@command -v air >/dev/null 2>&1 || go install github.com/air-verse/air@latest
	air

# Telnet client that waits for the server and reconnects after reloads
peer:
	@echo "peer → $(PEER_HOST):$(PEER_PORT) (auto-reconnect; Ctrl+C to quit)"
	@trap 'echo; echo "peer stopped."; exit 0' INT TERM; \
	while true; do \
		until nc -z $(PEER_HOST) $(PEER_PORT) 2>/dev/null; do sleep 0.3; done; \
		echo "connected to $(PEER_HOST):$(PEER_PORT)"; \
		telnet $(PEER_HOST) $(PEER_PORT) || true; \
		echo "disconnected — reconnecting…"; \
		sleep 0.5; \
	done

test:
	go test ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

tidy:
	go mod tidy

clean:
	rm -rf bin/ tmp/

help:
	@echo "Targets:"
	@echo "  build  - compile the server binary to $(BIN)"
	@echo "  run    - build and run the server"
	@echo "  dev    - run server with Air hot-reload"
	@echo "  peer   - telnet to $(PEER_HOST):$(PEER_PORT) (auto-reconnect)"
	@echo "  test   - run all tests"
	@echo "  fmt    - format Go source files"
	@echo "  vet    - run go vet"
	@echo "  tidy   - tidy go.mod and go.sum"
	@echo "  clean  - remove build artifacts"
	@echo "  help   - show this help message"
