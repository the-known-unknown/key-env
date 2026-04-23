BINARY  := key-env
SRC     := ./cmd/keyenv
GOFLAGS ?=

.PHONY: build clean

build:
	go build $(GOFLAGS) -o $(BINARY) $(SRC)

clean:
	rm -f $(BINARY)
