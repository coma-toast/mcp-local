BINARY := mcp-local
BINDIR := bin
BIN := $(BINDIR)/$(BINARY)
PREFIX ?= /usr/local
INSTALL_DIR ?= $(PREFIX)/bin
# Version stamped into the binary so `mcp-local update` can compare it with releases:
# "0.1.8" on a release tag, "0.1.8-3-gabc1234[-dirty]" after it, "dev" outside git.
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null | sed 's/^v//')
LDFLAGS := -X main.version=$(or $(VERSION),dev)

# GOWORK=off: this repo is not part of any go.work workspace, but Go's workspace
# auto-detection walks up from cwd, so a go.work in a parent directory (e.g. ~/git/go.work
# listing unrelated modules) would otherwise break module resolution.
export GOWORK := off

.PHONY: all build install clean test check fmt

all build:
	mkdir -p $(BINDIR)
	go build -ldflags "$(LDFLAGS)" -o $(BIN) ./cmd/mcp-local

install: build
	mkdir -p $(INSTALL_DIR)
	install -m 0755 $(BIN) $(INSTALL_DIR)/$(BINARY)

test:
	go test ./...

# gofmt -l exits non-zero on a parse error, which fails the target with gofmt's own message.
check: test
	go vet ./...
	@unformatted="$$(gofmt -l .)" || exit 1; if [ -n "$$unformatted" ]; then echo "unformatted files (run make fmt):"; echo "$$unformatted"; exit 1; fi

fmt:
	gofmt -w .

clean:
	rm -f $(BIN)
