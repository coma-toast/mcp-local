BINARY := mcp-local
BINDIR := bin
BIN := $(BINDIR)/$(BINARY)

.PHONY: all build install clean

all build:
	mkdir -p $(BINDIR)
	go build -o $(BIN) ./cmd/mcp-local

install: build
	install -m 0755 $(BIN) /usr/local/bin/

clean:
	rm -f $(BIN)
