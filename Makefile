BINARY  := key-env
SRC     := ./cmd/keyenv
GOFLAGS ?=

.PHONY: test build integration-test integration-test-py export-test clean release

# Run all Go unit tests
test:
	go test ./...

# Run unit tests, then compile the binary
build: test
	go build $(GOFLAGS) -o $(BINARY) $(SRC)

# Build, then verify the binary can resolve a KeePass secret end-to-end.
# Requires keepassxc-cli on PATH (used by the vault backend).
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

# Same as integration-test but does a clean build first.
# Does not require keepassxc-cli (uses the Go-native KeePass reader).
integration-test-py: clean build
	@expected="Jcg5TfdI9X0zHaU03Qx9bGb0rphYh0xIebtpFPTcRT"; \
	actual=$$(./$(BINARY) run \
		--env test/.env.sample \
		--secrets test/keepass-sample-db.kdbx \
		--password '4jFU%i*+Q2qdpFgoHJGK' \
		-- sh -c 'echo $$TEST_CLIENT_SECRET'); \
	if [ "$$actual" = "$$expected" ]; then \
		echo "integration-test-py: PASS"; \
	else \
		echo "integration-test-py: FAIL (expected '$$expected', got '$$actual')"; \
		exit 1; \
	fi

EXPORT_DIR := $(HOME)/.key-env-test

# Clean build, then bundle the binary + test fixtures into ~/.key-env-test
# so you can test key-env outside this repo. Includes a ready-to-run run.sh.
export-test: clean build
	@rm -rf $(EXPORT_DIR)
	@mkdir -p $(EXPORT_DIR)
	@cp -a test/. $(EXPORT_DIR)/
	@cp $(BINARY) $(EXPORT_DIR)/
	@printf '#!/bin/sh\ncd "$$( dirname "$$0" )"\n./key-env run \\\n    --env .env.sample \\\n    --secrets keepass-sample-db.kdbx \\\n    --password '\''4jFU%%i*+Q2qdpFgoHJGK'\'' \\\n    --verbose \\\n    -- node sample.js\n' > $(EXPORT_DIR)/run.sh
	@chmod +x $(EXPORT_DIR)/run.sh
	@echo "Exported to $(EXPORT_DIR):"
	@ls -1 $(EXPORT_DIR)

# Remove the compiled binary
clean:
	rm -f $(BINARY)

# Tag and push a release (usage: make release VERSION=0.1.0)
release:
	@test -n "$(VERSION)" || (echo "usage: make release VERSION=0.1.0" && exit 1)
	git tag -a v$(VERSION) -m "Release v$(VERSION)"
	git push origin v$(VERSION)
