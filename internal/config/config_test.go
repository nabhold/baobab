package config

import "testing"

func TestLoadRequiresSecureOIDCIssuer(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://example")
	t.Setenv("ADMIN_OIDC_AUDIENCE", "baobab-control-plane")
	t.Setenv("ADMIN_OIDC_ISSUER", "http://identity.example.com")
	t.Setenv("WORKLOAD_OIDC_AUDIENCE", "baobab-control-plane")
	t.Setenv("WORKLOAD_OIDC_ISSUER", "https://workload-identity.example.com")
	if _, err := Load(); err == nil {
		t.Fatal("insecure remote issuer was accepted")
	}
	t.Setenv("ADMIN_OIDC_ISSUER", "http://127.0.0.1:5556")
	if _, err := Load(); err != nil {
		t.Fatalf("local development issuer rejected: %v", err)
	}
	t.Setenv("WORKLOAD_OIDC_ISSUER", "http://workload-identity.example.com")
	if _, err := Load(); err == nil {
		t.Fatal("insecure workload issuer was accepted")
	}
}
