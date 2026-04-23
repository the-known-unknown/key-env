BINARY  := key-env
SRC     := ./cmd/keyenv
GOFLAGS ?=

.PHONY: test build integration-test clean release

test:
	go test ./...

build: test
	go build $(GOFLAGS) -o $(BINARY) $(SRC)

integration-test: build
	@command -v keepassxc-cli >/dev/null 2>&1 || (echo "error: keepassxc-cli not found on PATH" && exit 1)
	@expected="Jcg5TfdI9X0zHaU03Qx9bGb0rphYh0xIebtpFPTcRT"; \
	actual=$$(./$(BINARY) run \
		--env test/.env.sample \
		--secrets test/keepass-sample-db.kdbx \
		--password '4jFU%i*+Q2qdpFgoHJGK' \
		-- sh -c 'echo $$TEST_CLIENT_SECRET'); \
	if [ "$$actual" = "$$expected" ]; then \
		echo "integration-test: PASS"; \
	else \
		echo "integration-test: FAIL (expected '$$expected', got '$$actual')"; \
		exit 1; \
	fi

clean:
	rm -f $(BINARY)

release:
	@test -n "$(VERSION)" || (echo "usage: make release VERSION=0.1.0" && exit 1)
	git tag -a v$(VERSION) -m "Release v$(VERSION)"
	git push origin v$(VERSION)
