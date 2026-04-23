BINARY  := key-env
SRC     := ./cmd/keyenv
GOFLAGS ?=

.PHONY: build clean release

build:
	go build $(GOFLAGS) -o $(BINARY) $(SRC)

clean:
	rm -f $(BINARY)

release:
	@test -n "$(VERSION)" || (echo "usage: make release VERSION=0.1.0" && exit 1)
	git tag -a v$(VERSION) -m "Release v$(VERSION)"
	git push origin v$(VERSION)
