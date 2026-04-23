package keepassxc

import "testing"

func TestParseKPPath(t *testing.T) {
	secretPath, credential, err := parseKPPath("Services/DB/Main/Password")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if secretPath != "Services/DB/Main" {
		t.Fatalf("unexpected secret path: %q", secretPath)
	}
	if credential != "Password" {
		t.Fatalf("unexpected credential: %q", credential)
	}
}

func TestParseKPPathError(t *testing.T) {
	_, _, err := parseKPPath("Password")
	if err == nil {
		t.Fatal("expected error")
	}
}
