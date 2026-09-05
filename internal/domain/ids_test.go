package domain

import "testing"

func TestNewTenantIDMatchesSharedGrammar(t *testing.T) {
	for i := 0; i < 5; i++ {
		id := NewTenantID()
		if !ValidTenantID(id) {
			t.Fatalf("minted tenant id %q does not match the shared tn_ grammar", id)
		}
	}
}

func TestValidTenantID(t *testing.T) {
	cases := map[string]bool{
		"tn_01k4example": true,
		"tn_abc123":      true,
		"":               false,
		"tn_":            false,
		"tn_ABC123":      false, // uppercase not permitted after the prefix
		"tn_abc-123":     false, // hyphen not permitted
		"acme":           false, // missing prefix entirely
		"zuribeans_za":   false, // legacy-shaped alias, not a tenant id
	}
	for input, want := range cases {
		if got := ValidTenantID(input); got != want {
			t.Errorf("ValidTenantID(%q) = %v, want %v", input, got, want)
		}
	}
}

func TestValidLegalEntityID(t *testing.T) {
	cases := map[string]bool{
		"THAMANI-GLOBAL": true,
		"NABHOLD":        true,
		"zuribeans_za":   true, // accepted legacy alias per ADR-0003 §2.3
		"":               false,
		"thamani-global": false, // lowercase canonical form is not valid
		"THAMANI GLOBAL": false, // space not permitted
		"-THAMANI":       false, // must start with a letter
	}
	for input, want := range cases {
		if got := ValidLegalEntityID(input); got != want {
			t.Errorf("ValidLegalEntityID(%q) = %v, want %v", input, got, want)
		}
	}
}

func TestIsCanonicalLegalEntityID(t *testing.T) {
	if !IsCanonicalLegalEntityID("THAMANI-GLOBAL") {
		t.Fatal("canonical form not recognised as canonical")
	}
	if IsCanonicalLegalEntityID("zuribeans_za") {
		t.Fatal("legacy alias incorrectly recognised as canonical")
	}
}

func TestMintedOpaqueIDsMatchSharedGrammars(t *testing.T) {
	if id := NewMappingID(); !ValidMappingID(id) {
		t.Fatalf("minted mapping id %q does not match the shared map_ grammar", id)
	}
	if id := NewExternalReferenceID(); !ValidExternalReferenceID(id) {
		t.Fatalf("minted external reference id %q does not match the shared ref_ grammar", id)
	}
	if id := NewMappingScopeID(); !ValidMappingScopeID(id) {
		t.Fatalf("minted mapping scope id %q does not match the shared scope_ grammar", id)
	}
}

func TestValidProductID(t *testing.T) {
	cases := map[string]bool{
		"baobab-trade": true,
		"baobab_trade": true,
		"ab":           false, // below minimum length
		"Baobab-Trade": false, // uppercase not permitted
		"":             false,
	}
	for input, want := range cases {
		if got := ValidProductID(input); got != want {
			t.Errorf("ValidProductID(%q) = %v, want %v", input, got, want)
		}
	}
}
