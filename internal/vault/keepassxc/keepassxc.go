package keepassxc

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

type Provider struct {
	secretsPath string
	password    string
}

func New(secretsPath, password string) *Provider {
	return &Provider{
		secretsPath: secretsPath,
		password:    password,
	}
}

const cliName = "keepassxc-cli"

func (p *Provider) Resolve(path string) (string, error) {
	if strings.TrimSpace(path) == "" {
		return "", errors.New("empty KeepassXC path")
	}
	secretPath, credential, err := parseKPPath(path)
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(p.password) == "" {
		return "", errors.New("missing KeepassXC password; pass --password or use a secure password source")
	}

	if _, err := exec.LookPath(cliName); err != nil {
		return "", fmt.Errorf("%s is not installed.\n\n  Install it with:  brew install keepassxc\n\n  For other platforms, see: https://keepassxc.org/download", cliName)
	}

	cmd := exec.Command(cliName, "show", p.secretsPath, secretPath, "-q", "-s", "-a", credential)
	cmd.Stdin = strings.NewReader(p.password + "\n")

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = err.Error()
		}
		return "", fmt.Errorf("%s failed: %s", cliName, msg)
	}

	value := strings.TrimSpace(stdout.String())
	if value == "" {
		return "", fmt.Errorf("resolved empty value for path %q credential %q", secretPath, credential)
	}
	return value, nil
}

func parseKPPath(path string) (secretPath string, credential string, err error) {
	trimmed := strings.Trim(path, "/")
	if trimmed == "" {
		return "", "", errors.New("invalid kp path: empty")
	}
	last := strings.LastIndex(trimmed, "/")
	if last <= 0 || last >= len(trimmed)-1 {
		return "", "", errors.New("invalid kp path: expected kp://<secret_path>/<credential>")
	}
	return trimmed[:last], trimmed[last+1:], nil
}
