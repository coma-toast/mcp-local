BINARY := mcp-local

build:
	go build -o $(BINARY) ./cmd/mcp-local/main.go

install: build
	mv $(BINARY) /usr/local/bin/

clean:
	rm -f $(BINARY)
