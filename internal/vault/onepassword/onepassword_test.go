package onepassword

import "testing"

func TestParseOPPath(t *testing.T) {
	vault, item, field, err := parseOPPath("Shared/Stripe/api_key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if vault != "Shared" {
		t.Fatalf("unexpected vault: %q", vault)
	}
	if item != "Stripe" {
		t.Fatalf("unexpected item: %q", item)
	}
	if field != "api_key" {
		t.Fatalf("unexpected field: %q", field)
	}
}

func TestParseOPPathFieldWithSlash(t *testing.T) {
	vault, item, field, err := parseOPPath("Shared/db/sections/credentials/password")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if vault != "Shared" {
		t.Fatalf("unexpected vault: %q", vault)
	}
	if item != "db" {
		t.Fatalf("unexpected item: %q", item)
	}
	if field != "sections/credentials/password" {
		t.Fatalf("unexpected field: %q", field)
	}
}

func TestParseOPPathErrorEmpty(t *testing.T) {
	if _, _, _, err := parseOPPath(""); err == nil {
		t.Fatal("expected error")
	}
}

func TestParseOPPathErrorSingleSegment(t *testing.T) {
	if _, _, _, err := parseOPPath("Shared"); err == nil {
		t.Fatal("expected error")
	}
}

func TestParseOPPathErrorTwoSegments(t *testing.T) {
	if _, _, _, err := parseOPPath("Shared/Stripe"); err == nil {
		t.Fatal("expected error")
	}
}

func TestParseOPPathErrorEmptySegment(t *testing.T) {
	if _, _, _, err := parseOPPath("Shared//api_key"); err == nil {
		t.Fatal("expected error for empty middle segment")
	}
}
