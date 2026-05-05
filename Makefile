BINARY  := key-env
SRC     := ./cmd/keyenv
GOFLAGS ?=

.PHONY: test build integration-test integration-test-fail integrations integration-test-py export-test clean release

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
		--env test/.env.kp \
		--secrets test/keepass-sample-db.kdbx \
		--password '4jFU%i*+Q2qdpFgoHJGK' \
		-- sh -c 'echo $$TEST_CLIENT_SECRET'); \
	if [ "$$actual" = "$$expected" ]; then \
		echo "integration-test: PASS"; \
	else \
		echo "integration-test: FAIL (expected '$$expected', got '$$actual')"; \
		exit 1; \
	fi

# Verify that key-env surfaces child process failures correctly.
# The child exits with code 42; key-env should exit non-zero and
# its stderr should contain the original exit code.
integration-test-fail: build
	@command -v keepassxc-cli >/dev/null 2>&1 || (echo "error: keepassxc-cli not found on PATH" && exit 1)
	@stderr=$$(./$(BINARY) run \
		--env test/.env.kp \
		--secrets test/keepass-sample-db.kdbx \
		--password '4jFU%i*+Q2qdpFgoHJGK' \
		-- sh -c 'echo $$TEST_CLIENT_SECRET; exit 42' 2>&1 1>/dev/null); \
	rc=$$?; \
	if [ "$$rc" -eq 0 ]; then \
		echo "integration-test-fail: FAIL (expected non-zero exit, got 0)"; \
		exit 1; \
	fi; \
	case "$$stderr" in \
		*"exit status 42"*) echo "integration-test-fail: PASS (child exit 42 surfaced correctly)" ;; \
		*) echo "integration-test-fail: FAIL (stderr missing 'exit status 42', got: $$stderr)"; exit 1 ;; \
	esac

# Run all integration tests
integrations: integration-test integration-test-fail

# Same as integration-test but does a clean build first.
# Does not require keepassxc-cli (uses the Go-native KeePass reader).
integration-test-py: clean build
	@expected="Jcg5TfdI9X0zHaU03Qx9bGb0rphYh0xIebtpFPTcRT"; \
	actual=$$(./$(BINARY) run \
		--env test/.env.kp \
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
export-test:
	@rm -f $(BINARY)
	@go test ./... >/dev/null 2>&1
	@go build $(GOFLAGS) -o $(BINARY) $(SRC)
	@rm -rf $(EXPORT_DIR)
	@mkdir -p $(EXPORT_DIR)
	@cp -a test/. $(EXPORT_DIR)/
	@cp $(BINARY) $(EXPORT_DIR)/
	@printf '#!/bin/sh\ncd "$$( dirname "$$0" )"\n./key-env run \\\n    --env .env.kp \\\n    --secrets keepass-sample-db.kdbx \\\n    --password '\''4jFU%%i*+Q2qdpFgoHJGK'\'' \\\n    --verbose \\\n    -- node sample.js\n' > $(EXPORT_DIR)/run.sh
	@chmod +x $(EXPORT_DIR)/run.sh
	@echo "Exported to $(EXPORT_DIR)"

# Remove the compiled binary
clean:
	rm -f $(BINARY)

# Tag and push a release (usage: make release VERSION=0.1.0)
release:
	@test -n "$(VERSION)" || (echo "usage: make release VERSION=0.1.0" && exit 1)
	git tag -a v$(VERSION) -m "Release v$(VERSION)"
	git push origin v$(VERSION)
