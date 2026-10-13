BINARY := mcp-local
BINDIR := bin
BIN := $(BINDIR)/$(BINARY)
export GOWORK := off

.PHONY: all build install clean test check fmt

all build:
	mkdir -p $(BINDIR)
	GOWORK=$(GOWORK) go build -o $(BIN) ./cmd/mcp-local

install: build
	install -m 0755 $(BIN) $(INSTALL_DIR)/$(BINARY)

test:
	GOWORK=$(GOWORK) go test ./...

check: test
	go vet ./...
	@if [ -n "$$(gofmt -l . 2>/dev/null)" ]; then echo "Error: unformatted files:"; gofmt -l .; exit 1; fi

fmt:
	gofmt -w .

clean:
	rm -f $(BIN)

INSTALL_DIR ?= /usr/local/bin
