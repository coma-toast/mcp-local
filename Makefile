BINARY := mcp-local
BINDIR := bin
BIN := $(BINDIR)/$(BINARY)

.PHONY: all build install clean test check

all build:
	mkdir -p $(BINDIR)
	go build -o $(BIN) ./cmd/mcp-local

install: build
	install -m 0755 $(BIN) /usr/local/bin/

test:
	go test ./...

check: test
	go vet ./...
	go fmt ./...

clean:
	rm -f $(BIN)
