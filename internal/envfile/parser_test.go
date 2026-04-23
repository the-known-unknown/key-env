package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env.dev")
	content := `# comment
SECURE_SECRET_VAR="kp://path/to/entry/Password"
DB_NAME="userdb"
DB_PASSWORD="op://Vault/Path/To/Secret/Credential"
`
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	out, err := ParseFile(path)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if len(out) != 3 {
		t.Fatalf("expected 3 vars, got %d", len(out))
	}

	if out[0].Var != "SECURE_SECRET_VAR" || out[0].Type != TypeKP || out[0].Path != "path/to/entry/Password" {
		t.Fatalf("unexpected first var: %+v", out[0])
	}
	if out[1].Var != "DB_NAME" || out[1].Type != TypePlain || out[1].Path != "userdb" {
		t.Fatalf("unexpected second var: %+v", out[1])
	}
	if out[2].Var != "DB_PASSWORD" || out[2].Type != TypeOP || out[2].Path != "Vault/Path/To/Secret/Credential" {
		t.Fatalf("unexpected third var: %+v", out[2])
	}
}
