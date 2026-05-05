package onepassword

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

type Provider struct{}

func New() *Provider {
	return &Provider{}
}

const cliName = "op"

func (p *Provider) Resolve(path string) (string, error) {
	if strings.TrimSpace(path) == "" {
		return "", errors.New("empty 1Password path")
	}
	if _, _, _, err := parseOPPath(path); err != nil {
		return "", err
	}

	if _, err := exec.LookPath(cliName); err != nil {
		return "", fmt.Errorf("%s is not installed.\n\n  Install it with:  brew install 1password-cli\n\n  For other platforms, see: https://1password.com/downloads/command-line/", cliName)
	}

	cmd := exec.Command(cliName, "read", "op://"+path)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = err.Error()
		}
		if isUnauthenticated(msg) {
			return "", fmt.Errorf("%s is not authenticated: %s\n\n  Enable desktop app integration (1Password → Settings → Developer → \"Integrate with 1Password CLI\")\n  or set OP_SERVICE_ACCOUNT_TOKEN for headless use.\n  See: https://developer.1password.com/docs/cli/app-integration/", cliName, msg)
		}
		return "", fmt.Errorf("%s failed: %s", cliName, msg)
	}

	value := strings.TrimSpace(stdout.String())
	if value == "" {
		return "", fmt.Errorf("resolved empty value for op://%s", path)
	}
	return value, nil
}

func parseOPPath(path string) (vault string, item string, field string, err error) {
	trimmed := strings.Trim(path, "/")
	if trimmed == "" {
		return "", "", "", errors.New("invalid op path: empty")
	}
	parts := strings.SplitN(trimmed, "/", 3)
	if len(parts) < 3 {
		return "", "", "", errors.New("invalid op path: expected op://<vault>/<item>/<field>")
	}
	if parts[0] == "" || parts[1] == "" || parts[2] == "" {
		return "", "", "", errors.New("invalid op path: empty segment")
	}
	return parts[0], parts[1], parts[2], nil
}

func isUnauthenticated(msg string) bool {
	lower := strings.ToLower(msg)
	hints := []string{
		"not signed in",
		"no account",
		"session expired",
		"authorization prompt dismissed",
	}
	for _, h := range hints {
		if strings.Contains(lower, h) {
			return true
		}
	}
	return false
}
