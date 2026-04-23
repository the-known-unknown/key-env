package onepassword

import "fmt"

type Provider struct{}

func New() *Provider {
	return &Provider{}
}

func (p *Provider) Resolve(path string) (string, error) {
	return "", fmt.Errorf("1Password (op://) support is not yet available.\n\n  Track progress at: https://github.com/the-known-unknown/key-env")
}
