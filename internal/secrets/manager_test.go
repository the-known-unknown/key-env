package secrets

import (
	"errors"
	"testing"

	"github.com/asimmittal/key-env/internal/envfile"
	"github.com/asimmittal/key-env/internal/vault"
)

type fakeProvider struct {
	value string
	err   error
}

func (f fakeProvider) Resolve(path string) (string, error) {
	if f.err != nil {
		return "", f.err
	}
	return f.value + ":" + path, nil
}

func TestLoaderLoadSuccess(t *testing.T) {
	loader := NewLoader(map[string]vault.Provider{
		"kp": fakeProvider{value: "value"},
	}, false)
	parsed := []envfile.ParsedVar{
		{Var: "DB_NAME", Type: envfile.TypePlain, Path: "userdb"},
		{Var: "SECURE_SECRET_VAR", Type: envfile.TypeKP, Path: "path/to/secret/Password"},
	}

	out, err := loader.Load(parsed)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 values, got %d", len(out))
	}
	if out[0].Value != "userdb" || out[1].Value != "value:path/to/secret/Password" {
		t.Fatalf("unexpected values: %+v", out)
	}
}

func TestLoaderLoadProviderError(t *testing.T) {
	loader := NewLoader(map[string]vault.Provider{
		"kp": fakeProvider{err: errors.New("missing path")},
	}, false)
	parsed := []envfile.ParsedVar{
		{Var: "SECURE_SECRET_VAR", Type: envfile.TypeKP, Path: "path/to/secret/Password"},
	}

	_, err := loader.Load(parsed)
	if err == nil {
		t.Fatal("expected error")
	}
}
