package onepassword

import "fmt"

type Provider struct{}

func New() *Provider {
	return &Provider{}
}

func (p *Provider) Resolve(path string) (string, error) {
	return "", fmt.Errorf("onepassword provider not implemented yet for path %q", path)
}
